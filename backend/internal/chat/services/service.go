package services

import (
	"errors"
	"sort"

	"zone/internal/chat/models"
	"zone/internal/websocket"
)

// defaultHistoryBatch is the page size used when no limit is given.
const defaultHistoryBatch = 10

// Repository is the data boundary the chat service depends on.
type Repository interface {
	InsertMessage(senderID int, msgType string, refID int, content string) (models.ChatMessage, error)
	CanSendDM(senderID, recipientID int) bool
	IsGroupMember(userID, groupID int) bool
	EverMessaged(userID, otherID int) bool
	ListDM(userID, otherID, beforeID, limit int) ([]models.ChatMessage, bool, error)
	ListGroup(groupID, beforeID, limit int) ([]models.ChatMessage, bool, error)
	FollowGraphUsers(userID int) ([]models.ConversationUser, error)
	GroupMemberships(userID int) ([]models.ConversationGroup, error)
	LastMessageFor(userID, targetID int, msgType string) (*models.LastMessage, error)
}

// Pusher delivers chat envelopes over the realtime backbone.
type Pusher interface {
	SendToUser(userID int, env websocket.Envelope)
	SendToRoom(roomID int, env websocket.Envelope)
	IsOnline(userID int) bool
}

// Service holds business rules for chat; it has no HTTP and no SQL.
type Service struct {
	repo Repository
	push Pusher
}

// NewService builds a chat service with its repository and pusher.
func NewService(repo Repository, push Pusher) *Service {
	return &Service{repo: repo, push: push}
}
