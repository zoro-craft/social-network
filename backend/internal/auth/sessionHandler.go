package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

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
	w.Header().Set("Content-Type", "application/json")

	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "error",
			"message": "Unauthorized. Please login.",
		})
		return
	}

	var username, email, firstName, lastName, avatar string
	err := database.Database.QueryRow(
		"SELECT COALESCE(nickname, ''), email, COALESCE(first_name, ''), COALESCE(last_name, ''), COALESCE(avatar, '') FROM users WHERE id = ?",
		userID,
	).Scan(&username, &email, &firstName, &lastName, &avatar)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "error",
			"message": "User not found or database error.",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
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