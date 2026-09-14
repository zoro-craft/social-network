package repositories

import (
	"zone/internal/database"
	"zone/internal/groups/models"
)

func SaveGroup(groupe models.Groupe) (int, error) {
	result, err := database.Database.Exec(
		`INSERT INTO user_groups (creator_id, title, description) VALUES (?, ?, ?)`,
		groupe.Admin, groupe.Title, groupe.Description,
	)
	if err != nil {
		return 0, err
	}
	groupID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(groupID), err
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


func GetAllGroups(userID int) ([]models.Groupe, error) {
	rows, err := database.Database.Query(`
		SELECT
			ug.id,
			ug.creator_id,
			ug.title,
			ug.description,
			EXISTS (
				SELECT 1
				FROM group_members gm
				WHERE gm.group_id = ug.id
				AND gm.user_id = ?
			) AS is_member
		FROM user_groups ug
		ORDER BY ug.id DESC
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
			&group.IsMember,
		)

		if err != nil {
			return nil, err
		}

		groups = append(groups, group)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return groups, nil
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


func GetUserFullName(user_id int) string {
	var fullName string

	err := database.Database.QueryRow(`
		SELECT full_name
		FROM users
		WHERE id = ?
	`, user_id).Scan(&fullName)
	if err != nil {
		return ""
	}

	return fullName
}



func GetGroupbyid(group_id int) (models.Groupe, error) {
	var group models.Groupe

	err := database.Database.QueryRow(`
		SELECT id, creator_id, title, description
		FROM user_groups
		WHERE id = ?
	`, group_id).Scan(
		&group.Id,
		&group.Admin,
		&group.Title,
		&group.Description,
	)

	if err != nil {
		return group, err
	}

	return group, nil
}


func IsMember(userID int, groupID int) bool {
	var exists bool

	err := database.Database.QueryRow(`
		SELECT EXISTS (
			SELECT 1
			FROM group_members
			WHERE user_id = ? AND group_id = ?
		)
	`, userID, groupID).Scan(&exists)

	if err != nil {
		return false
	}

	return exists
}