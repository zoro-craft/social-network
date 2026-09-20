package auth

import (
	"log"
	"net/http"
	"zone/internal/utils"

	"zone/internal/database"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	cookie, err := r.Cookie("session_token")

	// delete session from database
	if err == nil && cookie.Value != "" {
		
		tx, txErr := database.Database.Begin()
		if txErr != nil {
			log.Println("Error starting transaction:", txErr)
		} else {
			defer tx.Rollback()

			query := "DELETE FROM user_sessions WHERE session_token = ?"
			_, err := tx.Exec(query, cookie.Value)

			if err != nil {
				log.Println("Error deleting session from DB:", err)
			} else {
				if commitErr := tx.Commit(); commitErr != nil {
					log.Println("Error committing transaction:", commitErr)
				}
			}
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

	utils.RespondJSON(w, http.StatusOK, map[string]string{"status": "success", "message": "Logged out successfully"})
}