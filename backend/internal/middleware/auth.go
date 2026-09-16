package middleware

import (
	"context"
	"net/http"
	"time"
	"zone/internal/utils"

	"zone/internal/database"
)

type contextKey string

const UserIDKey contextKey = "userID"

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			utils.HandleError(w, http.StatusUnauthorized, "Unauthorized - No session cookie")
			return
		}

		var userID int
		var expiresAt time.Time

		err = database.Database.QueryRow("SELECT user_id, expires_at FROM user_sessions WHERE session_token = ?", cookie.Value).Scan(&userID, &expiresAt)
		if err != nil || time.Now().After(expiresAt) {
			http.SetCookie(w, &http.Cookie{
				Name:     "session_token",
				Value:    "",
				Expires:  time.Now().Add(-1 * time.Hour),
				HttpOnly: true,
				Path:     "/",
			})
			utils.HandleError(w, http.StatusUnauthorized, "Unauthorized - Invalid or expired session")
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}