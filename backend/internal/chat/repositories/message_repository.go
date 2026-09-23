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

// IsGroupMember reports whether the user belongs to the group.
func (r ChatRepository) IsGroupMember(userID, groupID int) bool {
	var exists bool
	err := database.Database.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM group_members WHERE user_id = ? AND group_id = ?
		)`,
		userID, groupID,
	).Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}

// EverMessaged reports whether the two users have ever exchanged a private message.
func (r ChatRepository) EverMessaged(userID, otherID int) bool {
	var exists bool
	err := database.Database.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM messages
			WHERE type = 'private'
			  AND ((sender_id = ? AND ref_id = ?) OR (sender_id = ? AND ref_id = ?))
		)`,
		userID, otherID, otherID, userID,
	).Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}
