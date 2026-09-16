package repositories

import (
	"database/sql"
	"fmt"

	"zone/internal/database"
	"zone/internal/groups/models"
)

// insert the join request  in the group_requests table
func SaveJoinRequest(groupID int, userID int) (int, error) {
	var id int

	// Check if request already exists
	err := database.Database.QueryRow(
		`
		SELECT id
		FROM group_requests
		WHERE group_id = ?
		AND user_id = ?
	`,
		groupID,
		userID,
	).Scan(&id)

	if err == nil {
		return 0, fmt.Errorf("join request already exists")
	}

	if err != sql.ErrNoRows {
		return 0, err
	}

	// Insert request
	result, err := database.Database.Exec(
		`
		INSERT INTO group_requests (group_id, user_id)
		VALUES (?, ?)
	`,
		groupID,
		userID,
	)
	if err != nil {
		return 0, err
	}

	id64, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id64), nil
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
