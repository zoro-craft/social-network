package profile

import "zone/internal/posts"

type Profile struct {
	ID          int    `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	DateOfBirth string `json:"date_of_birth"`

	Nickname string `json:"nickname,omitempty"`
	AboutMe  string `json:"about_me,omitempty"`
	Avatar   string `json:"avatar,omitempty"`

	IsPublic bool `json:"is_public"`

	Posts []posts.Post `json:"posts"`
}
