package reactions

type Reaction struct {
	PostID       int    `json:"post_id"`
	Type         string `json:"type,omitempty"`
	LikeCount    int    `json:"like_count,omitempty"`
	DislikeCount int    `json:"dislike_count,omitempty"`
	UserReaction string `json:"user_reaction,omitempty"`
}