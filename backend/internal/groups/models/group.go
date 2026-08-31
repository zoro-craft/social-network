package models

type Group struct {
	ID          int    `json:"id"`
	Creator     int    `json:"creator_ID"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at,omitempty"`

	CreatorName string `json:"creator_name,omitempty"`
	MemberCount int    `json:"member_count,omitempty"`
	IsMember    bool   `json:"is_member,omitempty"`
	IsCreator   bool   `json:"is_creator,omitempty"`
}

type GroupMember struct {
	ID       int    `json:"id"`
	GroupID  int    `json:"group_id"`
	UserID   int    `json:"user_id"`
	JoinedAt string `json:"joined_at,omitempty"`

	FirstName string `json:"first_name,omitempty"`
	Avatar    string `json:"avatar,omitempty"`
	Nickname  string `json:"nickname,omitempty"`
}

type GroupInvitation struct {
	ID        int    `json:"id"`
	GroupID   int    `json:"group_id"`
	InviterID int    `json:"inviter_id"`
	UserID    int    `json:"invited_user_id"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at,omitempty"`

	GroupTitle  string `json:"group_title,omitempty"`
	GroupDesc   string `json:"group_description,omitempty"`
	CreatorName string `json:"creator_name,omitempty"`
}

type GroupRequest struct {
	ID        int    `json:"id"`
	GroupID   int    `json:"group_id"`
	UserID    int    `json:"user_id"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at,omitempty"`

	FirstName string `json:"first_name,omitempty"`
	Nickname  string `json:"nickname,omitempty"`
	Avatar    string `json:"avatar,omitempty"`
}
