package handlers

import (
	"net/http"

	"github.com/gorilla/websocket"

	"zone/internal/auth"
	"zone/internal/utils"
	ws "zone/internal/websocket"
	"zone/internal/websocket/repositories"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// WSHandler upgrades an authenticated request to a websocket connection.
func WSHandler(hub *ws.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := auth.GetUserID(r)
		if err != nil {
			utils.RespondJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
			return
		}

		// reject the upgrade before the socket opens when at the per-user cap
		if !hub.CanAccept(userID) {
			utils.RespondJSON(w, http.StatusTooManyRequests, map[string]string{"error": "Too many connections"})
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}

		client := ws.NewClient(hub, conn, userID)
		hub.Register(client)

		// subscribe to all the user's groups so group chat reaches them
		if groupIDs, err := repositories.GetUserGroupIDs(userID); err == nil {
			for _, gid := range groupIDs {
				hub.Subscribe(userID, gid)
			}
		}

		go client.WritePump()
		client.ReadPump()
	}
}