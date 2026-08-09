package auth

import (
	"net/http"

	"zone/internal/database"
	"zone/internal/utils"
)

func Me(w http.ResponseWriter, r *http.Request) {
	userID, err := GetUserIDFromSession(r)
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	var firstName string
	var nickname, avatar *string
	err = database.Database.QueryRow(
		"SELECT first_name, nickname, avatar FROM users WHERE id = ?", userID,
	).Scan(&firstName, &nickname, &avatar)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Server error")
		return
	}

	nicknameValue := ""
	if nickname != nil {
		nicknameValue = *nickname
	}

	avatarValue := ""
	if avatar != nil {
		avatarValue = *avatar
	}

	utils.RespondJSON(w, http.StatusOK, map[string]any{
		"authenticated": true,
		"id":            userID,
		"first_name":    firstName,
		"nickname":      nicknameValue,
		"avatar":        avatarValue,
	})
}