package websocket

// hubMessage is the union of control messages the hub goroutine handles.
type hubMessage struct {
	kind   string
	client *Client
	userID int
	roomID int
	env    Envelope
	count  int
	resp   chan int
}

// Hub is the shared realtime registry: one goroutine owns all state.
type Hub struct {
	connections  map[int]map[*Client]bool
	rooms        map[int]map[*Client]bool
	register     chan hubMessage
	unregister   chan hubMessage
	broadcast    chan hubMessage
	query        chan hubMessage
	done         chan struct{}
	onConnect    func(userID int)
	onDisconnect func(userID int)
}

// DefaultHub is the process-wide hub singleton, mirroring database.Database.
var DefaultHub = NewHub()

// NewHub builds an empty hub.
func NewHub() *Hub {
	return &Hub{
		connections:  make(map[int]map[*Client]bool),
		rooms:        make(map[int]map[*Client]bool),
		register:     make(chan hubMessage),
		unregister:   make(chan hubMessage),
		broadcast:    make(chan hubMessage),
		query:        make(chan hubMessage),
		done:         make(chan struct{}),
		onConnect:    func(int) {},
		onDisconnect: func(int) {},
	}
}

// Run is the single goroutine that owns the hub's maps.
func (h *Hub) Run() {
	for {
		select {
		case msg := <-h.register:
			h.handleRegister(msg)
		case msg := <-h.unregister:
			h.handleUnregister(msg)
		case msg := <-h.broadcast:
			switch msg.kind {
			case "user":
				h.handleSendUser(msg)
			case "room":
				h.handleSendRoom(msg)
			case "subscribe":
				h.handleSubscribe(msg)
			}
		case msg := <-h.query:
			switch msg.kind {
			case "online":
				msg.resp <- btoi(h.isOnline(msg.userID))
			case "count":
				msg.resp <- len(h.connections[msg.userID])
			}
		case <-h.done:
			h.handleCloseAll()
			return
		}
	}
}

// handleCloseAll closes every connection; read pumps exit and the hub returns.
func (h *Hub) handleCloseAll() {
	for _, set := range h.connections {
		for c := range set {
			c.conn.Close()
		}
	}
}

// btoi converts a bool to 1/0 for the query channel.
func btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}

func (h *Hub) handleRegister(msg hubMessage) {
	set := h.connections[msg.client.userID]
	if len(set) >= maxConnectionsPerUser {
		return
	}
	wasOffline := len(set) == 0
	if set == nil {
		set = make(map[*Client]bool)
		h.connections[msg.client.userID] = set
	}
	set[msg.client] = true
	if wasOffline {
		h.onConnect(msg.client.userID)
	}
}

func (h *Hub) handleUnregister(msg hubMessage) {
	// remove the dying socket from every room, dropping empty rooms
	for roomID, set := range h.rooms {
		delete(set, msg.client)
		if len(set) == 0 {
			delete(h.rooms, roomID)
		}
	}

	set := h.connections[msg.client.userID]
	if set == nil {
		close(msg.client.send)
		return
	}
	delete(set, msg.client)
	if len(set) == 0 {
		delete(h.connections, msg.client.userID)
		h.onDisconnect(msg.client.userID)
	}
	// close the send channel so WritePump exits via !ok; the client is already out of every set.
	close(msg.client.send)
}

// sendToSet delivers an envelope to every client in a set, never blocking.
func (h *Hub) sendToSet(set map[*Client]bool, env Envelope) {
	for c := range set {
		select {
		case c.send <- env:
		default:
		}
	}
}

func (h *Hub) handleSendUser(msg hubMessage) {
	h.sendToSet(h.connections[msg.userID], msg.env)
}

func (h *Hub) handleSendRoom(msg hubMessage) {
	h.sendToSet(h.rooms[msg.roomID], msg.env)
}

func (h *Hub) handleSubscribe(msg hubMessage) {
	if h.rooms[msg.roomID] == nil {
		h.rooms[msg.roomID] = make(map[*Client]bool)
	}
	for c := range h.connections[msg.userID] {
		h.rooms[msg.roomID][c] = true
	}
}

func (h *Hub) isOnline(userID int) bool {
	return len(h.connections[userID]) > 0
}

// Register adds a client to the hub.
func (h *Hub) Register(c *Client) {
	h.sendIfOpen(h.register, hubMessage{kind: "register", client: c})
}

// Unregister removes a client from the hub.
func (h *Hub) Unregister(c *Client) {
	h.sendIfOpen(h.unregister, hubMessage{kind: "unregister", client: c})
}

// SendToUser pushes an envelope to every connection of a user.
func (h *Hub) SendToUser(userID int, env Envelope) {
	h.sendIfOpen(h.broadcast, hubMessage{kind: "user", userID: userID, env: env})
}

// SendToRoom pushes an envelope to every connection in a room.
func (h *Hub) SendToRoom(roomID int, env Envelope) {
	h.sendIfOpen(h.broadcast, hubMessage{kind: "room", roomID: roomID, env: env})
}

// Subscribe registers a user's connections for a room.
func (h *Hub) Subscribe(userID, roomID int) {
	h.sendIfOpen(h.broadcast, hubMessage{kind: "subscribe", userID: userID, roomID: roomID})
}

// SetPresenceHooks installs callbacks fired on first-connect and last-disconnect.
func (h *Hub) SetPresenceHooks(onConnect, onDisconnect func(userID int)) {
	h.onConnect = onConnect
	h.onDisconnect = onDisconnect
}

// IsOnline reports whether a user has at least one open connection.
func (h *Hub) IsOnline(userID int) bool {
	return h.queryBool("online", userID) == 1
}

// ConnectionCount reports the number of open connections for a user.
func (h *Hub) ConnectionCount(userID int) int {
	return h.queryBool("count", userID)
}

// CanAccept reports whether the user has room for another connection.
func (h *Hub) CanAccept(userID int) bool {
	return h.ConnectionCount(userID) < maxConnectionsPerUser
}

// queryBool sends a query message and waits for the hub's reply.
func (h *Hub) queryBool(kind string, userID int) int {
	resp := make(chan int, 1)
	select {
	case h.query <- hubMessage{kind: kind, userID: userID, resp: resp}:
		return <-resp
	case <-h.done:
		return 0
	}
}

// sendIfOpen delivers a message unless the hub is shutting down.
func (h *Hub) sendIfOpen(dst chan hubMessage, msg hubMessage) {
	select {
	case dst <- msg:
	case <-h.done:
	}
}

// Close shuts the hub down, closing all connections and draining.
func (h *Hub) Close() {
	close(h.done)
}
