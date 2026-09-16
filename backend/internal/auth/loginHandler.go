package auth

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
	"zone/internal/utils"

	"zone/internal/database"

	"golang.org/x/crypto/bcrypt"
)

const (
	SessionDuration = 24 * time.Hour
)

type LoginReq struct {
	Identify string `json:"identify"`
	Password string `json:"password"`
}

// creates a random 32 byte session token.
func generateSessionToken() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		log.Println("Error generating session token:", err)
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var identify, password string
	contentType := r.Header.Get("Content-Type")

	if strings.HasPrefix(contentType, "application/json") {
		var req LoginReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.HandleError(w, http.StatusBadRequest, "Invalid request formatting")
			return
		}
		identify = req.Identify
		password = req.Password
	} else {
		if err := r.ParseForm(); err != nil {
			utils.HandleError(w, http.StatusBadRequest, "Invalid form data")
			return
		}
		identify = r.FormValue("login")
		if identify == "" {
			identify = r.FormValue("identify")
		}
		password = r.FormValue("password")
	}

	var dbID int
	var dbUsername string
	var dbPasswordHash string
	var discardEmail string

	err := database.Database.QueryRow(`
		SELECT id, password, COALESCE(nickname, ''), email FROM users 
		WHERE email = ? OR (nickname = ? AND nickname != '')
	`, identify, identify).Scan(&dbID, &dbPasswordHash, &dbUsername, &discardEmail)

	if err != nil {
		if err == sql.ErrNoRows {
			utils.HandleError(w, http.StatusUnauthorized, "Invalid credentials")
			return
		}
		utils.HandleError(w, http.StatusInternalServerError, "Server error")
		return
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(dbPasswordHash), []byte(password))
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Clear old sessions
	deleteQuery := "DELETE FROM user_sessions WHERE user_id = ?"
	_, err = database.Database.Exec(deleteQuery, dbID)
	if err != nil {
		log.Println("Error deleting old sessions:", err)
	}

	// Create new session token
	sessionToken := generateSessionToken()
	if sessionToken == "" {
		utils.HandleError(w, http.StatusInternalServerError, "Server error")
		return
	}
	expirationTime := time.Now().Add(SessionDuration)
	_, err = database.Database.Exec("INSERT INTO user_sessions (user_id, session_token, expires_at) VALUES (?, ?, ?)", dbID, sessionToken, expirationTime)
	if err != nil {
		log.Println("Error saving session to DB:", err)
		utils.HandleError(w, http.StatusInternalServerError, "Error creating session, try again later")
		return
	}

	// Set session cookie
	newCookie := http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  expirationTime,
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &newCookie)

	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "Logged in successfully",
		"user": map[string]interface{}{
			"id":       dbID,
			"nickname": dbUsername,
		},
	})
}
