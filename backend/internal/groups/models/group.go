package models

import "time"

type Groupe struct {
	Id          int    `json:"id"`
	Admin       int    `json:"admin"`
	Title       string `json:"titel"`
	Description string `json:"description"`
	IsMember    bool   `json:"is_member"`
}

type GroupInvitation struct {
	ID            int       `json:"id"`
	GroupID       int       `json:"group_id"`
	InviterID     int       `json:"inviter_id"`
	InvitedUserID int       `json:"invited_user_id"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}


type GroupRequest struct {
	ID        int       `json:"id"`
	GroupID   int       `json:"group_id"`
	UserID    int       `json:"user_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type User struct {
	ID        int    `json:"id"`
	Full_name string `json:"full_name"`
	Avatar    string `json:"avatar,omitempty"`
}
