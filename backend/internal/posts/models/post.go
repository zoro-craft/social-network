package models

type Post struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	Content      string `json:"content"`
	Image        string `json:"image"`
	Privacy      string `json:"privacy"`
	GroupID      *int   `json:"group_id,omitempty"`
	UserID       int    `json:"user_id"`
	FullName     string `json:"full_name"`
	Avatar       string `json:"avatar"`
	CreatedAt    string `json:"created_at"`
	LikeCount    int    `json:"like_count"`
	DislikeCount int    `json:"dislike_count"`
	UserReaction string `json:"user_reaction,omitempty"`
}
