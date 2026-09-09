package utils

import (
	"zone/internal/database"
)

func IsGroupMember(groupID, userID int) bool {
	var isMember bool
	err := database.Database.QueryRow(
		"SELECT EXISTS(SELECT * FROM group_members WHERE group_id = ? AND user_id = ?)",
		groupID, userID,
	).Scan(&isMember)
	return err == nil && isMember
}
