package websocket

import (
	"time"

	"github.com/gorilla/websocket"
)

const (
	// maxConnectionsPerUser caps concurrent sockets per userID.
	maxConnectionsPerUser = 10
	// writeWait is the time allowed to write a message to the peer.
	writeWait = 10 * time.Second
	// pongWait is the time allowed to read the next pong from the peer.
	pongWait = 10 * time.Second
	// pingPeriod is how often the server pings the peer.
	pingPeriod = 9 * time.Second
	// sendBuffer is the outbound queue capacity per connection.
	sendBuffer = 64
)

// Client is a single websocket connection with its pumps.
type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	userID int
	send   chan Envelope
}

// NewClient builds a Client for a user.
func NewClient(h *Hub, conn *websocket.Conn, userID int) *Client {
	return &Client{hub: h, conn: conn, userID: userID, send: make(chan Envelope, sendBuffer)}
}

// ReadPump reads messages until the connection closes or the pong deadline passes.
func (c *Client) ReadPump() {
	defer c.hub.Unregister(c)
	defer c.conn.Close()

	c.conn.SetReadLimit(4096)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		return c.conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			return
		}
	}
}

// WritePump writes queued envelopes and periodic pings to the socket.
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case env, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			data, err := MarshalEnvelope(env)
			if err != nil {
				continue
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}