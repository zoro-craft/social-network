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

// ConversationUser is a people row in the conversation sidebar.
type ConversationUser struct {
	UserID        int     `json:"user_id"`
	Name          string  `json:"name"`
	Avatar        string  `json:"avatar"`
	Online        bool    `json:"online"`
	LastMessage   *string `json:"last_message"`
	LastMessageAt *int64  `json:"last_message_at"`
	CanChat       bool    `json:"can_chat"`
}

// ConversationGroup is a group row in the conversation sidebar.
type ConversationGroup struct {
	GroupID       int     `json:"group_id"`
	Title         string  `json:"title"`
	LastMessage   *string `json:"last_message"`
	LastMessageAt *int64  `json:"last_message_at"`
}

// LastMessage is the newest message preview for a conversation.
type LastMessage struct {
	Content   string
	CreatedAt int64
}