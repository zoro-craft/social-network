package repositories

import (
	"zone/internal/database"
	"zone/internal/groups/models"
)

// insert the join request  in the group_requests table

func SaveJoinRequest(group_id int, user_id int) (int, error) {
	result, err := database.Database.Exec(`
		INSERT INTO group_requests (group_id, user_id)
		VALUES (?, ?)
	`, group_id, user_id)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func GetJoinRequest(requestID int) (models.GroupRequest, error) {
	var request models.GroupRequest

	err := database.Database.QueryRow(`
		SELECT id, group_id, user_id, status, created_at
		FROM group_requests
		WHERE id = ?
	`, requestID).Scan(
		&request.ID,
		&request.GroupID,
		&request.UserID,
		&request.Status,
		&request.CreatedAt,
	)

	return request, err
}


func UpdateJoinRequestStatus(requestID int, status string) error {
	_, err := database.Database.Exec(`
		UPDATE group_requests
		SET status = ?
		WHERE id = ?
	`, status, requestID)
	if err != nil {
		return err
	}

	return nil
}