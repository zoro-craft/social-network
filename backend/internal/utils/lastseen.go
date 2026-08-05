package utils

import (
	"log"

	"zone/internal/database"
)

func UpdateLastSeen(userID int) {
	_, err := database.Database.Exec(
		"UPDATE users SET last_seen = CURRENT_TIMESTAMP WHERE id = ?",
		userID,
	)
	if err != nil {
		log.Println("UpdateLastSeen:", err)
	}
}