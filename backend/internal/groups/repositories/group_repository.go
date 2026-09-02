package repositories

import (
	"zone/internal/database"
	"zone/internal/groups/models"
)

func SaveGroup(groupe models.Groupe) error {
	_, err := database.Database.Exec(
		`INSERT INTO user_groups (creator_id, title, description) VALUES (?, ?, ?)`,
		groupe.Admin, groupe.Titel, groupe.Description,
	)
	return err
}


func GetGroupNameByID(groupID int) (string, error) {
    var groupName string

    err := database.Database.QueryRow(`
        SELECT title
        FROM user_groups
        WHERE id = ?
    `, groupID).Scan(&groupName)

    if err != nil {
        return "", err
    }

    return groupName, nil
}

func GroupExists(groupID int) bool {
	var exists bool

	err := database.Database.QueryRow(
		`SELECT EXISTS(
			SELECT 1 FROM user_groups WHERE id = ?
		)`,
		groupID,
	).Scan(&exists)
	if err != nil {
		return false
	}

	return exists
}

func AddGroupMember(groupID, userID int) error {
	_, err := database.Database.Exec(`
		INSERT INTO group_members (group_id, user_id)
		VALUES (?, ?)
	`, groupID, userID)

	return err
}