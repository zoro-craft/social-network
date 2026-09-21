package repositories

import (
	"zone/internal/database"
)

// MemberIDsByGroup returns every member user id of a group for the event fan-out.
func (r NotificationRepository) MemberIDsByGroup(groupID int) ([]int, error) {
	rows, err := database.Database.Query(
		"SELECT user_id FROM group_members WHERE group_id = ?",
		groupID,
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
