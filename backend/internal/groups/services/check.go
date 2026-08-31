package services

import (
	"strings"

	"zone/internal/database"
	"zone/internal/groups/models"
)

func ValidGroup(group models.Group) bool {
	if len(strings.TrimSpace(group.Title)) <= 0 || len(strings.TrimSpace(group.Title)) >= 50 {
		return false
	}
	if len(strings.TrimSpace(group.Description)) <= 0 || len(strings.TrimSpace(group.Description)) >= 200 {
		return false
	}
	return true
}

func ValidEvent(event models.Event) bool {
	if len(strings.TrimSpace(event.Title)) <= 0 || len(strings.TrimSpace(event.Title)) >= 100 {
		return false
	}
	if strings.TrimSpace(event.EventTime) == "" {
		return false
	}
	return true
}

func IsGroupMember(groupID, userID int) bool {
	var exists bool
	database.Database.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM group_members WHERE group_id = ? AND user_id = ?)",
		groupID, userID,
	).Scan(&exists)
	return exists
}

func IsGroupCreator(groupID, userID int) bool {
	var exists bool
	database.Database.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM user_groups WHERE id = ? AND creator_id = ?)",
		groupID, userID,
	).Scan(&exists)
	return exists
}

func GroupExists(groupID int) bool {
	var exists bool
	database.Database.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM user_groups WHERE id = ?)",
		groupID,
	).Scan(&exists)
	return exists
}
