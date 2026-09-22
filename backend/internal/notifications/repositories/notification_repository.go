package repositories

import (
	"database/sql"
	"time"

	"zone/internal/database"
	"zone/internal/notifications/models"
)

// NotificationRepository is the only layer that reads and writes the notifications table.
type NotificationRepository struct{}

// NewNotificationRepository builds a notification repository.
func NewNotificationRepository() NotificationRepository {
	return NotificationRepository{}
}

const listLimit = 10

// Insert persists a notification row with an app-set epoch-ms created_at and returns the row.
func (r NotificationRepository) Insert(recipientID int, notifType string, refID *int, actorID *int, message string) (models.Notification, error) {
	createdAt := time.Now().UnixMilli()
	res, err := database.Database.Exec(`
		INSERT INTO notifications (user_id, actor_id, type, ref_id, message, is_read, created_at)
		VALUES (?, ?, ?, ?, ?, 0, ?)`,
		recipientID, actorID, notifType, refID, message, createdAt,
	)
	if err != nil {
		return models.Notification{}, err
	}
	id, _ := res.LastInsertId()
	return models.Notification{
		ID:        int(id),
		UserID:    recipientID,
		ActorID:   actorID,
		Type:      notifType,
		RefID:     refID,
		Message:   message,
		IsRead:    false,
		CreatedAt: createdAt,
	}, nil
}

const selectColumns = `
	SELECT n.id, n.user_id, n.actor_id, n.type, n.ref_id, n.message, n.is_read, n.created_at
	FROM notifications n`

// InsertFollowRequest persists a follow_request notification for (recipient, requester),
// reusing an existing row (reset to unread, refreshed) so repeated sends never duplicate.
func (r NotificationRepository) InsertFollowRequest(recipientID, requesterID int, message string) (models.Notification, error) {
	createdAt := time.Now().UnixMilli()
	var id int
	err := database.Database.QueryRow(`
		SELECT id FROM notifications
		WHERE user_id = ? AND type = 'follow_request' AND ref_id = ?
		LIMIT 1`, recipientID, requesterID).Scan(&id)

	switch {
	case err == sql.ErrNoRows:
		res, err := database.Database.Exec(`
			INSERT INTO notifications (user_id, actor_id, type, ref_id, message, is_read, created_at)
			VALUES (?, ?, 'follow_request', ?, ?, 0, ?)`,
			recipientID, requesterID, requesterID, message, createdAt)
		if err != nil {
			return models.Notification{}, err
		}
		inserted, _ := res.LastInsertId()
		id = int(inserted)
	case err != nil:
		return models.Notification{}, err
	default:
		_, err = database.Database.Exec(`
			UPDATE notifications SET is_read = 0, message = ?, actor_id = ?, created_at = ?
			WHERE id = ?`,
			message, requesterID, createdAt, id)
		if err != nil {
			return models.Notification{}, err
		}
	}

	var row models.Notification
	var actorID, refID sql.NullInt64
	err = database.Database.QueryRow(selectColumns+" WHERE n.id = ?", id).
		Scan(&row.ID, &row.UserID, &actorID, &row.Type, &refID, &row.Message, &row.IsRead, &row.CreatedAt)
	if err != nil {
		return models.Notification{}, err
	}
	if actorID.Valid {
		a := int(actorID.Int64)
		row.ActorID = &a
	}
	if refID.Valid {
		r := int(refID.Int64)
		row.RefID = &r
	}
	return row, nil
}

// ListUnread returns up to 10 unread notifications for a user, newest first.
// Keyset cursor on id: cursorID=0 means the first page. hasMore reports whether
// more unread rows exist beyond the returned page.
func (r NotificationRepository) ListUnread(userID int, cursorID int) ([]models.Notification, bool, error) {
	query := selectColumns + " WHERE n.user_id = ? AND n.is_read = 0"
	args := []any{userID}
	if cursorID > 0 {
		query += " AND n.id < ?"
		args = append(args, cursorID)
	}
	query += " ORDER BY n.id DESC LIMIT ?"
	args = append(args, listLimit+1)

	rows, err := database.Database.Query(query, args...)
	if err != nil {
		return nil, false, err
	}
	defer rows.Close()

	items := []models.Notification{}
	for rows.Next() {
		var n models.Notification
		var actorID sql.NullInt64
		var refID sql.NullInt64
		if err := rows.Scan(&n.ID, &n.UserID, &actorID, &n.Type, &refID, &n.Message, &n.IsRead, &n.CreatedAt); err != nil {
			return nil, false, err
		}
		if actorID.Valid {
			id := int(actorID.Int64)
			n.ActorID = &id
		}
		if refID.Valid {
			id := int(refID.Int64)
			n.RefID = &id
		}
		items = append(items, n)
	}
	if err := rows.Err(); err != nil {
		return nil, false, err
	}

	hasMore := len(items) > listLimit
	if hasMore {
		items = items[:listLimit]
	}
	return items, hasMore, nil
}

// CountUnread returns the exact total unread count for a user, independent of pagination.
func (r NotificationRepository) CountUnread(userID int) (int, error) {
	var count int
	err := database.Database.QueryRow(
		"SELECT COUNT(*) FROM notifications WHERE user_id = ? AND is_read = 0",
		userID,
	).Scan(&count)
	return count, err
}

// MarkReadOnAction marks a notification read by (user, type, ref_id) after accept/decline.
func (r NotificationRepository) MarkReadOnAction(userID int, notifType string, refID int) error {
	_, err := database.Database.Exec(`
		UPDATE notifications SET is_read = 1
		WHERE user_id = ? AND type = ? AND ref_id = ?`,
		userID, notifType, refID,
	)
	return err
}

// MarkReadByID marks a single notification read, owner-scoped; ErrNoRows when not owned.
func (r NotificationRepository) MarkReadByID(userID, id int) error {
	res, err := database.Database.Exec(
		"UPDATE notifications SET is_read = 1 WHERE id = ? AND user_id = ?",
		id, userID,
	)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// DeleteOnAction deletes a notification by (user, type, ref_id) and returns its id,
// or 0 when no row matches.
func (r NotificationRepository) DeleteOnAction(userID int, notifType string, refID int) (int, error) {
	var id int
	err := database.Database.QueryRow(
		"SELECT id FROM notifications WHERE user_id = ? AND type = ? AND ref_id = ?",
		userID, notifType, refID,
	).Scan(&id)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	if _, err := database.Database.Exec("DELETE FROM notifications WHERE id = ?", id); err != nil {
		return 0, err
	}
	return id, nil
}
