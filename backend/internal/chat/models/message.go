package models

// Message is a persisted chat message row.
type Message struct {
	ID        int    `json:"id"`
	SenderID  int    `json:"sender_id"`
	Type      string `json:"type"`
	RefID     int    `json:"ref_id"`
	Content   string `json:"content"`
	CreatedAt int64  `json:"created_at"`
}

// ChatMessage is the realtime payload shared by chat_new_message and chat_new_group_message.
type ChatMessage struct {
	ID         int    `json:"id"`
	SenderID   int    `json:"sender_id"`
	SenderName string `json:"sender_name"`
	Content    string `json:"content"`
	CreatedAt  int64  `json:"created_at"`
	Type       string `json:"type"`
	RefID      int    `json:"ref_id"`
}
