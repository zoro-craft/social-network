package stories

// Story is a single ephemeral image post that disappears after 24 hours.
type Story struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	Image     string `json:"image"`
	Caption   string `json:"caption"`
	CreatedAt int64  `json:"created_at"`
	Seen      bool   `json:"seen"`
}

// StoryGroup buckets one author's active stories, mirroring the Instagram rail
// where a single avatar opens that author's whole sequence.
type StoryGroup struct {
	UserID   int     `json:"user_id"`
	FullName string  `json:"full_name"`
	Avatar   string  `json:"avatar"`
	Seen     bool    `json:"seen"`
	Stories  []Story `json:"stories"`
}
