package posts

type Post struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	Content      string `json:"content"`
	Image        string `json:"image"`
	Privacy      string `json:"privacy"`
	UserID       int    `json:"user_id"`
	Full_name     string `json:"full_name"`
	CreatedAt    string `json:"created_at"`
	LikeCount    int    `json:"like_count"`
	DislikeCount int    `json:"dislike_count"`
	UserReaction string `json:"user_reaction,omitempty"` 
}