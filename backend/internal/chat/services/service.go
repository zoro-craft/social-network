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

// SendDM validates, persists, and delivers a direct message to the recipient and the sender.
func (s *Service) SendDM(userID, recipientID int, content string) error {
	if err := validateContent(content); err != nil {
		return err
	}
	if err := validateSelfDM(userID, recipientID); err != nil {
		return err
	}
	if !s.repo.CanSendDM(userID, recipientID) {
		return ErrNotAllowed
	}
	msg, err := s.repo.InsertMessage(userID, "private", recipientID, content)
	if err != nil {
		return err
	}
	env := websocket.Envelope{Type: "chat_new_message", Payload: msg}
	s.push.SendToUser(recipientID, env)
	s.push.SendToUser(userID, env)
	return nil
}

// SendGroup validates, persists, and delivers a group message to every connected member.
func (s *Service) SendGroup(userID, groupID int, content string) error {
	if err := validateContent(content); err != nil {
		return err
	}
	if !s.repo.IsGroupMember(userID, groupID) {
		return ErrNotMember
	}
	msg, err := s.repo.InsertMessage(userID, "group", groupID, content)
	if err != nil {
		return err
	}
	env := websocket.Envelope{Type: "chat_new_group_message", Payload: msg}
	s.push.SendToRoom(groupID, env)
	return nil
}

// ListHistory returns a newest-first batch of messages with has_more, gated by access.
func (s *Service) ListHistory(userID, targetID int, msgType string, beforeID, limit int) ([]models.ChatMessage, bool, error) {
	if limit <= 0 {
		limit = defaultHistoryBatch
	}
	switch msgType {
	case "private":
		if !s.repo.CanSendDM(userID, targetID) && !s.repo.EverMessaged(userID, targetID) {
			return nil, false, ErrNotAllowed
		}
		return s.repo.ListDM(userID, targetID, beforeID, limit)
	case "group":
		if !s.repo.IsGroupMember(userID, targetID) {
			return nil, false, ErrNotMember
		}
		return s.repo.ListGroup(targetID, beforeID, limit)
	default:
		return nil, false, errors.New("unknown message type")
	}
}

// GetConversations assembles the users and groups sidebar with online status,
// last-message previews, and ordering by last message desc then alphabetical.
func (s *Service) GetConversations(userID int) ([]models.ConversationUser, []models.ConversationGroup, error) {
	users, err := s.repo.FollowGraphUsers(userID)
	if err != nil {
		return nil, nil, err
	}
	groups, err := s.repo.GroupMemberships(userID)
	if err != nil {
		return nil, nil, err
	}

	for i := range users {
		users[i].Online = s.push.IsOnline(users[i].UserID)
		if lm, err := s.repo.LastMessageFor(userID, users[i].UserID, "private"); err == nil && lm != nil {
			content := lm.Content
			at := lm.CreatedAt
			users[i].LastMessage = &content
			users[i].LastMessageAt = &at
		}
	}
	for i := range groups {
		if lm, err := s.repo.LastMessageFor(userID, groups[i].GroupID, "group"); err == nil && lm != nil {
			content := lm.Content
			at := lm.CreatedAt
			groups[i].LastMessage = &content
			groups[i].LastMessageAt = &at
		}
	}

	sortUsers(users)
	sortGroups(groups)
	return users, groups, nil
}

// sortUsers orders people by last-message desc, then name asc.
func sortUsers(users []models.ConversationUser) {
	sort.SliceStable(users, func(i, j int) bool {
		ai, aj := users[i].LastMessageAt, users[j].LastMessageAt
		if ai != nil && aj != nil && *ai != *aj {
			return *ai > *aj
		}
		if (ai == nil) != (aj == nil) {
			return ai != nil
		}
		return users[i].Name < users[j].Name
	})
}

// sortGroups orders groups by last-message desc, then title asc.
func sortGroups(groups []models.ConversationGroup) {
	sort.SliceStable(groups, func(i, j int) bool {
		ai, aj := groups[i].LastMessageAt, groups[j].LastMessageAt
		if ai != nil && aj != nil && *ai != *aj {
			return *ai > *aj
		}
		if (ai == nil) != (aj == nil) {
			return ai != nil
		}
		return groups[i].Title < groups[j].Title
	})
}