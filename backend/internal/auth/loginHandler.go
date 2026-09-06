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

type apiResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
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

	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(apiResponse{Status: "error", Message: "Method not allowed"})
		return
	}

	var identify, password string
	contentType := r.Header.Get("Content-Type")

	if strings.HasPrefix(contentType, "application/json") {
		var req LoginReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(apiResponse{Status: "error", Message: "Invalid request formatting"})
			return
		}
		identify = req.Identify
		password = req.Password
	} else {
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(apiResponse{Status: "error", Message: "Invalid form data"})
			return
		}
		identify = r.FormValue("login")
		if identify == "" {
			identify = r.FormValue("identify")
		}
		password = r.FormValue("password")
	}

	var dbID int
	var dbUsername, dbEmail string
	var dbPasswordHash string

	err := database.Database.QueryRow(`
		SELECT id, password, COALESCE(nickname, ''), email FROM users 
		WHERE email = ? OR (nickname = ? AND nickname != '')
	`, identify, identify).Scan(&dbID, &dbPasswordHash, &dbUsername, &dbEmail)

	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(apiResponse{Status: "error", Message: "Invalid email/nickname or password"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(apiResponse{Status: "error", Message: "Server error"})
		return
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(dbPasswordHash), []byte(password))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(apiResponse{Status: "error", Message: "Invalid username/email or password."})
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
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(apiResponse{Status: "error", Message: "Server error"})
		return
	}
	expirationTime := time.Now().Add(SessionDuration)
	_, err = database.Database.Exec("INSERT INTO user_sessions (user_id, session_token, expires_at) VALUES (?, ?, ?)", dbID, sessionToken, expirationTime)
	if err != nil {
		log.Println("Error saving session to DB:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(apiResponse{Status: "error", Message: "Error creating session, try again later"})
		return
	}

	// Set session cookie
	newCookie := http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  expirationTime,
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(w, &newCookie)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Logged in successfully",
		"user": map[string]interface{}{
			"id":       dbID,
			"username": dbUsername,
			"nickname": dbUsername,
		},
	})
}
