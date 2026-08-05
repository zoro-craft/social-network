package auth

import (
	"net/http"
	"strings"
	"time"
	"zone/internal/database"
	"zone/internal/utils"

	"github.com/gofrs/uuid"

	"golang.org/x/crypto/bcrypt"
)
func Login(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login" {
		utils.HandleError(w, http.StatusNotFound, "Page not found")
		return
	}

	if r.Method != http.MethodPost {
		utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	identifier := strings.TrimSpace(r.FormValue("login"))
	password := r.FormValue("password")

	if identifier == "" || password == "" {
		utils.HandleError(w, http.StatusBadRequest, "Email/username and password are required")
		return
	}

	if len(identifier) < 2 || len(identifier) > 50 || len(password) < 6 || len(password) > 21 {
		utils.HandleError(w, http.StatusBadRequest, "Invalid data form")
		return
	}

	var userID int
	var hashedPassword string
	err := database.Database.QueryRow(
		"SELECT id, password FROM users WHERE email = ? OR nickname = ?",
		identifier, identifier,
	).Scan(&userID, &hashedPassword)
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, "Invalid email/username or password")
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		utils.HandleError(w, http.StatusUnauthorized, "Invalid email/username or password")
		return
	}

	_, err = database.Database.Exec("DELETE FROM sessions WHERE user_id = ?", userID)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Server error")
		return
	}

	sessionUUID, err := uuid.NewV4()
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Server error")
		return
	}
	sessionID := sessionUUID.String()
	expiration := time.Now().Add(24 * time.Hour)

	_, err = database.Database.Exec(
		"INSERT INTO sessions (id, expires_at, user_id) VALUES (?, ?, ?)",
		sessionID, expiration, userID,
	)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Server error")
		return
	}

	utils.UpdateLastSeen(userID)

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionID,
		Path:     "/",
		Expires:  expiration,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "login successfully"})
}