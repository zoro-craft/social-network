package repositories

import "zone/internal/database"

// GetUserGroupIDs returns the group IDs a user belongs to.
func GetUserGroupIDs(userID int) ([]int, error) {
	rows, err := database.Database.Query(
		"SELECT group_id FROM group_members WHERE user_id = ?",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}