package comments

type Comment struct {
	ID        int    `json:"id"`
	PostID    int    `json:"post_id"`
	GroupID   *int   `json:"group_id,omitempty"`
	UserID    int    `json:"user_id"`
	Content   string `json:"content"`
	Image     string `json:"image,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	Avatar    string `json:"avatar,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
}
