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

// ListDM returns a newest-first batch of private messages between two users.
func (r ChatRepository) ListDM(userID, otherID, beforeID, limit int) ([]models.ChatMessage, bool, error) {
	query := `
		SELECT m.id, m.sender_id, u.full_name, m.content, m.created_at, m.type, m.ref_id
		FROM messages m
		JOIN users u ON u.id = m.sender_id
		WHERE m.type = 'private'
		  AND ((m.sender_id = ? AND m.ref_id = ?) OR (m.sender_id = ? AND m.ref_id = ?))`
	args := []any{userID, otherID, otherID, userID}
	if beforeID > 0 {
		query += " AND m.created_at < (SELECT created_at FROM messages WHERE id = ?)"
		args = append(args, beforeID)
	}
	query += " ORDER BY m.created_at DESC LIMIT ?"
	args = append(args, limit+1)
	return queryMessages(limit, query, args...)
}

// ListGroup returns a newest-first batch of messages in a group.
func (r ChatRepository) ListGroup(groupID, beforeID, limit int) ([]models.ChatMessage, bool, error) {
	query := `
		SELECT m.id, m.sender_id, u.full_name, m.content, m.created_at, m.type, m.ref_id
		FROM messages m
		JOIN users u ON u.id = m.sender_id
		WHERE m.type = 'group' AND m.ref_id = ?`
	args := []any{groupID}
	if beforeID > 0 {
		query += " AND m.created_at < (SELECT created_at FROM messages WHERE id = ?)"
		args = append(args, beforeID)
	}
	query += " ORDER BY m.created_at DESC LIMIT ?"
	args = append(args, limit+1)
	return queryMessages(limit, query, args...)
}

// queryMessages runs a paged history query (limit+1 fetched) and returns rows plus has_more.
func queryMessages(limit int, query string, args ...any) ([]models.ChatMessage, bool, error) {
	rows, err := database.Database.Query(query, args...)
	if err != nil {
		return nil, false, err
	}
	defer rows.Close()

	items := []models.ChatMessage{}
	for rows.Next() {
		var m models.ChatMessage
		if err := rows.Scan(&m.ID, &m.SenderID, &m.SenderName, &m.Content, &m.CreatedAt, &m.Type, &m.RefID); err != nil {
			return nil, false, err
		}
		items = append(items, m)
	}
	if err := rows.Err(); err != nil {
		return nil, false, err
	}

	hasMore := len(items) > limit
	if hasMore {
		items = items[:limit]
	}
	return items, hasMore, nil
}