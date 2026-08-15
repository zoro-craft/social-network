package profile

import "zone/internal/posts"

type Profile struct {
	ID          int    `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email,omitempty"`
	DateOfBirth string `json:"date_of_birth,omitempty"`

	Nickname string `json:"nickname,omitempty"`
	AboutMe  string `json:"about_me,omitempty"`
	Avatar   string `json:"avatar,omitempty"`

	IsPublic   bool `json:"is_public"`
	Restricted bool `json:"restricted,omitempty"`

	FollowersCount int `json:"followers_count"`
	FollowingCount int `json:"following_count"`

	Posts []posts.Post `json:"posts"`
}