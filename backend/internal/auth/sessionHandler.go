package auth

import (
	"fmt"
	"net/http"
	"zone/internal/utils"

	"zone/internal/database"
	"zone/internal/middleware"
)

func GetUserID(r *http.Request) (int, error) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		return 0, fmt.Errorf("user ID not found in context")
	}
	return userID, nil
}

func MeHandler(w http.ResponseWriter, r *http.Request) {

	userID, err := GetUserID(r)
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, "Unauthorized. Please login.")
		return
	}

	var username, email, firstName, lastName, avatar string
	err = database.Database.QueryRow(
		"SELECT COALESCE(nickname, ''), email, COALESCE(first_name, ''), COALESCE(last_name, ''), COALESCE(avatar, '') FROM users WHERE id = ?",
		userID,
	).Scan(&username, &email, &firstName, &lastName, &avatar)
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, "User not found or database error.")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"status": "success",
		"user": map[string]interface{}{
			"id":         userID,
			"username":   username,
			"nickname":   username,
			"email":      email,
			"first_name": firstName,
			"last_name":  lastName,
			"avatar":     avatar,
		},
	})
}
