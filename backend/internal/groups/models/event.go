package models

type Event struct {
	ID          int    `json:"id"`
	GroupID     int    `json:"group_id"`
	UserID      int    `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	EventTime   string `json:"event_time"`
	CreatedAt   string `json:"created_at,omitempty"`

	CreatorName string `json:"creator_name,omitempty"`
	Going       int    `json:"going,omitempty"`
	NotGoing    int    `json:"not_going,omitempty"`
	Maybe       int    `json:"maybe,omitempty"`

	MyResponse string `json:"my_response,omitempty"`
}

type EventResponse struct {
	ID       int    `json:"id"`
	EventID  int    `json:"event_id"`
	UserID   int    `json:"user_id"`
	Response string `json:"response"`
}
