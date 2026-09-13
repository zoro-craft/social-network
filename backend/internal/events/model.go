package events

type VoteCount struct {
	Option string `json:"option"`
	Count  int    `json:"count"`
}

type Event struct {
	ID           int         `json:"id"`
	GroupID      int         `json:"group_id"`
	Title        string      `json:"title"`
	Description  string      `json:"description"`
	CreatedBy    int         `json:"created_by"`
	CreatedAt    string      `json:"created_at"`
	Options      []string    `json:"options"`
	VoteCounts   []VoteCount `json:"vote_counts"`
	UserResponse string      `json:"user_response,omitempty"`
}

type voteRequest struct {
	EventID int    `json:"event_id"`
	Vote    string `json:"vote"`
}

type createEventRequest struct {
	GroupID       int    `json:"group_id"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	FitsOptions   string `json:"fitsOptions"`
	SecondOptions string `json:"secondOptions"`
}