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
