package models

import "time"

type User struct {
	ID          int       `json:"id"`
	Email       string    `json:"email"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	DateOfBirth string    `json:"date_of_birth,omitempty"`
	Nickname    string    `json:"nickname,omitempty"`
	Avatar      string    `json:"avatar,omitempty"`
	AboutMe     string    `json:"about_me,omitempty"`
	IsPrivate   bool      `json:"is_private"`
	CreatedAt   time.Time `json:"created_at"`
}

type UserProfile struct {
	User      User   `json:"user"`
	Followers int    `json:"followers_count"`
	Following int    `json:"following_count"`
	Relation  string `json:"relation"` // "none", "pending", "accepted", "self"
}

type Follow struct {
	FollowerID int    `json:"follower_id"`
	FolloweeID int    `json:"followee_id"`
	Status     string `json:"status"`

	FirstName string `json:"first_name,omitempty"`
	Email     string `json:"email,omitempty"`
	Avatar    string `json:"avatar,omitempty"`
}
