package repositories

import (
	"zone/internal/database"
	"zone/internal/groups/models"
)

func GetNotificationsByUserID(userID int) ([]models.Notification, error) {
	rows, err := database.Database.Query(`
        SELECT id, user_id, type, invitation_id, message, is_read, created_at
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
