package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"social-network/pkg/db"
)

func GetUserID(r *http.Request) (int, error) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return 0, fmt.Errorf("no cookie found")
	}

	var userID int
	var expiresAt time.Time

	err = db.DB.QueryRow("SELECT user_id, expires_at FROM user_sessions WHERE session_token = ?", cookie.Value).Scan(&userID, &expiresAt)
	if err != nil {
		return 0, fmt.Errorf("invalid session token")
	}

	if time.Now().After(expiresAt) {
		return 0, fmt.Errorf("session expired")
	}

	return userID, nil
}

func SessionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, err := GetUserID(r)
	if err != nil {
		// Clear the invalid session cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
		})

		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "error",
			"message": "Unauthorized. Please login.",
		})
		return
	}

	var username, email, firstName, lastName string
	err = db.DB.QueryRow(
		"SELECT username, email, COALESCE(first_name, ''), COALESCE(last_name, '') FROM users WHERE id = ?",
		userID,
	).Scan(&username, &email, &firstName, &lastName)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "error",
			"message": "User not found.",
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
		},
	})
}