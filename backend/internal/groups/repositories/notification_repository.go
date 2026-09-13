package repositories

import (
	"zone/internal/database"
	"zone/internal/groups/models"
)


func CreateNotification(userID int, Type string, invitationID *int, requestID *int, msg string) error {
	_, err := database.Database.Exec(`
		INSERT INTO notifications
			(user_id, type, invitation_id, request_id, message)
		VALUES (?, ?, ?, ?, ?)
	`,
		userID,
		Type,
		invitationID,
		requestID,
		msg,
	)

	return err
}


func GetNotificationsByUserID(userID int) ([]models.Notification, error) {
	rows, err := database.Database.Query(`
        SELECT id, user_id, type,  COALESCE(invitation_id, 0) , COALESCE(request_id, 0),  message, is_read, created_at
        FROM notifications
        WHERE user_id = ?
        ORDER BY created_at DESC
    `, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []models.Notification

	for rows.Next() {
		var notification models.Notification

		err := rows.Scan(
			&notification.ID,
			&notification.UserID,
			&notification.Type,
			&notification.InvitationID,
			&notification.RequestID,
			&notification.Message,
			&notification.IsRead,
			&notification.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		notifications = append(notifications, notification)
	}

	return notifications, nil
}

func DeleteInvitationNotification(userID int, invitationID int) error {
	_, err := database.Database.Exec(
		`
		DELETE FROM notifications
		WHERE user_id = ?
		AND invitation_id = ?
		AND type = 'group_invtation'
	`,
		userID,
		invitationID,
	)

	return err
}

func DeleteRequestNotification(userID int, Requsestid int) error {
	_, err := database.Database.Exec(
		`
		DELETE FROM notifications
		WHERE user_id = ?
		AND  request_id = ?
		AND type = 'join_request'
	`,
		userID,
		Requsestid,
	)

	return err
}
