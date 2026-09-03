package repositories

import (
	"zone/internal/database"
	"zone/internal/groups/models"
)

func SaveGroup(groupe models.Groupe) error {
	_, err := database.Database.Exec(
		`INSERT INTO user_groups (creator_id, title, description) VALUES (?, ?, ?)`,
		groupe.Admin, groupe.Title, groupe.Description,
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


func GetAllGroupsNotMember(userID int) ([]models.Groupe, error) {
	rows, err := database.Database.Query(`
		SELECT 
			g.id,
			g.creator_id,
			g.title,
			g.description
		FROM user_groups g
		WHERE NOT EXISTS (
			SELECT 1
			FROM group_members gm
			WHERE gm.group_id = g.id
			AND gm.user_id = ?
		)
		ORDER BY g.id DESC
	`, userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []models.Groupe

	for rows.Next() {
		var group models.Groupe

		err := rows.Scan(
			&group.Id,
			&group.Admin,
			&group.Title,
			&group.Description,
		)
		if err != nil {
			return nil, err
		}

		groups = append(groups, group)
	}


	return groups, nil
}

func Alreadymember(group_id int, user_id int) bool {
	var exists bool

	err := database.Database.QueryRow(`
		SELECT id
		FROM group_members
		WHERE group_id = ? AND user_id = ?
	`, group_id, user_id).Scan(&exists)

	if err != nil {
		return false
	}

	return exists
}


func GetGroupCreator(group_id int) (int, error) {
	var creatorID int

	err := database.Database.QueryRow(`
		SELECT creator_id
		FROM user_groups
		WHERE id = ?
	`, group_id).Scan(&creatorID)

	if err != nil {
		return 0, err
	}

	return creatorID, nil
}