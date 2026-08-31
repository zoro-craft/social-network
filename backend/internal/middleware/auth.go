package middleware

import (
	"net/http"
	"time"
	"zone/internal/utils"
	"zone/internal/auth"
)

// Auth is a middleware that restricts access to authenticated users.
// It verifies the user's session token and updates their last seen activity.
func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := auth.GetUserIDFromSession(r)
		if err != nil {
			http.SetCookie(w, &http.Cookie{
				Name:     "session_token",
				Value:    "",
				Path:     "/",
				Expires:  time.Unix(0, 0),
				HttpOnly: true,
				SameSite: http.SameSiteLaxMode,
			})
			utils.HandleError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		// Update activity timestamp
		utils.UpdateLastSeen(userID)

		next(w, r)
	}
}