package websocket

import "encoding/json"

// Envelope is the single type-tagged wire message sent to clients.
type Envelope struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

// PresencePayload is the payload for type "presence".
type PresencePayload struct {
	UserID int  `json:"user_id"`
	Online bool `json:"online"`
}

// MarshalEnvelope encodes an envelope to JSON bytes for writing to a socket.
func MarshalEnvelope(env Envelope) ([]byte, error) {
	return json.Marshal(env)
}