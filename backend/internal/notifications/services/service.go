package services

import (
	"zone/internal/notifications/models"
	"zone/internal/notifications/repositories"
	"zone/internal/websocket"
)

// Notifier is the shared notification service wired to the real repository and realtime hub.
var Notifier = NewService(repositories.NewNotificationRepository(), websocket.DefaultHub)

// Repository is the data boundary the service depends on (implemented by repositories.NotificationRepository).
type Repository interface {
	Insert(recipientID int, notifType string, refID *int, actorID *int, message string) (models.Notification, error)
	InsertFollowRequest(recipientID int, requesterID int, message string) (models.Notification, error)
ListUnread(userID int, cursorID int) ([]models.Notification, bool, error)
	CountUnread(userID int) (int, error)
	MemberIDsByGroup(groupID int) ([]int, error)
	MarkReadOnAction(userID int, notifType string, refID int) error
	MarkReadByID(userID, id int) error
	DeleteOnAction(userID int, notifType string, refID int) (int, error)
}

// Pusher delivers an envelope to a user over the realtime backbone.
type Pusher interface {
	SendToUser(userID int, env websocket.Envelope)
}

// Service holds business rules for notifications; it has no HTTP and no SQL.
type Service struct {
	repo Repository
	push Pusher
}

// NewService builds a notification service with its repository and pusher.
func NewService(repo Repository, push Pusher) *Service {
	return &Service{repo: repo, push: push}
}

// Notify persists a notification for a recipient and pushes it over the websocket.
func (s *Service) Notify(recipientID int, notifType string, refID *int, actorID *int, message string) (models.Notification, error) {
	row, err := s.repo.Insert(recipientID, notifType, refID, actorID, message)
	if err != nil {
		return models.Notification{}, err
	}
	s.push.SendToUser(recipientID, websocket.Envelope{Type: "notification", Payload: row})
	return row, nil
}

// NotifyFollowRequest persists a follow_request notification for (recipient, requester) without duplicating, then pushes it live.
func (s *Service) NotifyFollowRequest(recipientID, requesterID int, message string) (models.Notification, error) {
	row, err := s.repo.InsertFollowRequest(recipientID, requesterID, message)
	if err != nil {
		return models.Notification{}, err
	}
	s.push.SendToUser(recipientID, websocket.Envelope{Type: "notification", Payload: row})
	return row, nil
}

// ListUnread returns the recipient's unread notifications, newest first, keyset-paginated on id.
func (s *Service) ListUnread(userID int, cursorID int) ([]models.Notification, bool, error) {
	return s.repo.ListUnread(userID, cursorID)
}

// CountUnread returns the exact total unread count for a user.
func (s *Service) CountUnread(userID int) (int, error) {
	return s.repo.CountUnread(userID)
}

// NotifyGroupEvent fans out an event_created notification to every group member except the creator.
func (s *Service) NotifyGroupEvent(groupID, eventID, creatorID int, title, creatorName string) error {
	members, err := s.repo.MemberIDsByGroup(groupID)
	if err != nil {
		return err
	}
	message := "New event by " + creatorName + " in your group: " + title
	for _, memberID := range members {
		if memberID == creatorID {
			continue
		}
		s.Notify(memberID, "event_created", &eventID, &creatorID, message)
	}
	return nil
}

// MarkReadOnAction marks a notification read when its entity is accepted or declined.
func (s *Service) MarkReadOnAction(userID int, notifType string, refID int) error {
	return s.repo.MarkReadOnAction(userID, notifType, refID)
}

// MarkReadByID marks a single notification read by id, owner-scoped.
func (s *Service) MarkReadByID(userID, id int) error {
	return s.repo.MarkReadByID(userID, id)
}

// RemoveNotificationOnAction deletes a notification by (recipient, type, ref) and signals the recipient.
func (s *Service) RemoveNotificationOnAction(recipientID int, notifType string, refID int) error {
	id, err := s.repo.DeleteOnAction(recipientID, notifType, refID)
	if err != nil {
		return err
	}
	if id == 0 {
		return nil
	}
	s.push.SendToUser(recipientID, websocket.Envelope{Type: "notification_removed", Payload: map[string]int{"id": id}})
	return nil
}

// RemoveFollowRequestNotification deletes the follow_request for (recipient, requester) and signals the recipient.
func (s *Service) RemoveFollowRequestNotification(followeeID, requesterID int) error {
	return s.RemoveNotificationOnAction(followeeID, "follow_request", requesterID)
}
