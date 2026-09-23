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

// InsertMessage persists a message with an app-set epoch-ms created_at and returns the full payload.
func (r ChatRepository) InsertMessage(senderID int, msgType string, refID int, content string) (models.ChatMessage, error) {
	createdAt := time.Now().UnixMilli()
	res, err := database.Database.Exec(`
		INSERT INTO messages (sender_id, type, ref_id, content, created_at)
		VALUES (?, ?, ?, ?, ?)`,
		senderID, msgType, refID, content, createdAt,
	)
	if err != nil {
		return models.ChatMessage{}, err
	}
	id, _ := res.LastInsertId()

	var row models.ChatMessage
	err = database.Database.QueryRow(`
		SELECT m.id, m.sender_id, u.full_name, m.content, m.created_at, m.type, m.ref_id
		FROM messages m
		JOIN users u ON u.id = m.sender_id
		WHERE m.id = ?`, id,
	).Scan(&row.ID, &row.SenderID, &row.SenderName, &row.Content, &row.CreatedAt, &row.Type, &row.RefID)
	if err != nil {
		return models.ChatMessage{}, err
	}
	return row, nil
}
