package models

type VoteCount struct {
	Option string `json:"option"`
	Count  int    `json:"count"`
}

type Event struct {
	ID           int         `json:"id"`
	GroupID      int         `json:"group_id"`
	Title        string      `json:"title"`
	Description  string      `json:"description"`
	EventDate    string      `json:"event_date"`
	EventTime    string      `json:"event_time"`
	CreatedBy    int         `json:"created_by"`
	CreatedAt    string      `json:"created_at"`
	Options      []string    `json:"options"`
	VoteCounts   []VoteCount `json:"vote_counts"`
	UserResponse string      `json:"user_response,omitempty"`
}

type VoteRequest struct {
	EventID int    `json:"event_id"`
	Vote    string `json:"vote"`
}

type CreateEventRequest struct {
	GroupID       int    `json:"group_id"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	EventDate     string `json:"eventDate"`
	EventTime     string `json:"eventTime"`
	FitsOptions   string `json:"fitsOptions"`
	SecondOptions string `json:"secondOptions"`
}
