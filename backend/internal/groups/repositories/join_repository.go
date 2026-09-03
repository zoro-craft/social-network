package repositories

import (
	"zone/internal/database"
)

// insert the join request  in the group_requests table

func SaveJoinRequest(group_id int, user_id int) error {
	_, err := database.Database.Exec(`
		INSERT INTO group_requests (group_id, user_id)
		VALUES (?, ?)
	`, group_id, user_id)

	return err
}
