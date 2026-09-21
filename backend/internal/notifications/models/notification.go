package models

// Notification is a persisted notification row.
type Notification struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	ActorID   *int   `json:"actor_id"`
	Type      string `json:"type"`
	RefID     *int   `json:"ref_id"`
	Message   string `json:"message"`
	IsRead    bool   `json:"is_read"`
	CreatedAt int64  `json:"created_at"`
}
