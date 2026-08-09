package posts

type Post struct {
	ID       int    `json:"id"`
	UserID   int    `json:"user_id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	Image    string `json:"image,omitempty"`
	Privacy  string `json:"privacy"`
	Nickname string `json:"nickname,omitempty"`

	CreatedAt string `json:"created_at,omitempty"`
}