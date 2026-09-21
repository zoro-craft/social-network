package websocket

import (
	"zone/internal/utils"
	"zone/internal/websocket/repositories"
)

// notifyPresence broadcasts presence{online:bool} to the user's audience.
func notifyPresence(h *Hub, userID int, online bool) {
	utils.UpdateLastSeen(userID)
	env := Envelope{
		Type:    "presence",
		Payload: PresencePayload{UserID: userID, Online: online},
	}
	for _, viewerID := range repositories.GetPresenceAudience(userID) {
		h.SendToUser(viewerID, env)
	}
}

// notifyConnect broadcasts presence{online:true} to the user's audience.
func notifyConnect(h *Hub, userID int) {
	notifyPresence(h, userID, true)
}

// notifyDisconnect broadcasts presence{online:false} to the user's audience.
func notifyDisconnect(h *Hub, userID int) {
	notifyPresence(h, userID, false)
}

// WirePresence installs presence hooks that run in their own goroutine.
func WirePresence(h *Hub) {
	h.SetPresenceHooks(
		func(userID int) { go notifyConnect(h, userID) },
		func(userID int) { go notifyDisconnect(h, userID) },
	)
}
