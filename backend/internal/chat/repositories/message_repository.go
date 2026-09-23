package repositories

import (
	"time"

	"zone/internal/chat/models"
	"zone/internal/database"
)

// CanSendDM reports whether an accepted follow exists in either direction.
func (r ChatRepository) CanSendDM(senderID, recipientID int) bool {
	var exists bool
	err := database.Database.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM followers
			WHERE (follower_id = ? AND followee_id = ? AND status = 'accepted')
			   OR (follower_id = ? AND followee_id = ? AND status = 'accepted')
		)`,
		senderID, recipientID, recipientID, senderID,
	).Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}
