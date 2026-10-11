package reels

// Reel is a single short video post. Unlike stories, reels do not expire: they
// make up a browsable vertical feed, so every engagement counter is persisted.
type Reel struct {
	ID           int    `json:"id"`
	UserID       int    `json:"user_id"`
	FullName     string `json:"full_name"`
	Avatar       string `json:"avatar"`
	Video        string `json:"video"`
	Caption      string `json:"caption"`
	CreatedAt    int64  `json:"created_at"`
	LikeCount    int    `json:"like_count"`
	CommentCount int    `json:"comment_count"`
	ViewCount    int    `json:"view_count"`
	SaveCount    int    `json:"save_count"`
	ShareCount   int    `json:"share_count"`
	Liked        bool   `json:"liked"`
	Saved        bool   `json:"saved"`
}

// ReelStats is the counter bundle returned after any engagement action, so the
// client can patch a single reel in place without refetching the whole feed.
type ReelStats struct {
	ReelID       int  `json:"reel_id"`
	LikeCount    int  `json:"like_count"`
	CommentCount int  `json:"comment_count"`
	ViewCount    int  `json:"view_count"`
	SaveCount    int  `json:"save_count"`
	ShareCount   int  `json:"share_count"`
	Liked        bool `json:"liked"`
	Saved        bool `json:"saved"`
}

// ReelComment is a text reply attached to a reel.
type ReelComment struct {
	ID        int    `json:"id"`
	ReelID    int    `json:"reel_id"`
	UserID    int    `json:"user_id"`
	Content   string `json:"content"`
	FirstName string `json:"first_name,omitempty"`
	Avatar    string `json:"avatar,omitempty"`
	CreatedAt int64  `json:"created_at"`
}
