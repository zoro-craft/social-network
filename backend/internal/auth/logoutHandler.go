package auth

import (
	"encoding/json"
	"log"
	"net/http"

	"zone/internal/database"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(apiResponse{Status: "error", Message: "Method not allowed"})
		return
	}

	cookie, err := r.Cookie("session_token")

	// delete session from database
	if err == nil && cookie.Value != "" {

		query := "DELETE FROM user_sessions WHERE session_token = ?"
		_, err := database.Database.Exec(query, cookie.Value)

		if err != nil {
			log.Println("Error deleting session from DB:", err)
		}
	}

	// expire session cookie
	deletedCookie := &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}

	http.SetCookie(w, deletedCookie)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(apiResponse{
		Status:  "success",
		Message: "Logged out successfully",
	})
}
