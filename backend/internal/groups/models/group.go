package models

import "time"

type Groupe struct {
	Id          int    `json:"id"`
	Admin       int    `json:"admin"`
	Titel       string `json:"titel"`
	Description string `json:"description"`
}

type GroupInvitation struct {
	ID            int       `json:"id"`
	GroupID       int       `json:"group_id"`
	InviterID     int       `json:"inviter_id"`
	InvitedUserID int       `json:"invited_user_id"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}

type Notification struct {
	ID           int       `json:"id"`
	UserID       int       `json:"user_id"`
	Type         string    `json:"type"`
	InvitationID int      `json:"invitation_id"`
	Message      string    `json:"message"`
	IsRead       bool      `json:"is_read"`
	CreatedAt    time.Time `json:"created_at"`
}