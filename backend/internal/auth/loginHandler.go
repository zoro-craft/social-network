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

	"zone/internal/auth/repositories"
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

// creates random session token
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

	tx, err := database.Database.Begin()
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Server error")
		return
	}
	defer tx.Rollback()

	creds, err := repositories.CredentialsByIdentifier(tx, identify)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.HandleError(w, http.StatusUnauthorized, "Invalid credentials")
			return
		}
		utils.HandleError(w, http.StatusInternalServerError, "Server error")
		return
	}

	// verify password
	err = bcrypt.CompareHashAndPassword([]byte(creds.Password), []byte(password))
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// clear old sessions
	if err := repositories.DeleteUserSessions(tx, creds.ID); err != nil {
		log.Println("Error deleting old sessions:", err)
	}

	// create new session token
	sessionToken := generateSessionToken()
	if sessionToken == "" {
		utils.HandleError(w, http.StatusInternalServerError, "Server error")
		return
	}
	expirationTime := time.Now().Add(SessionDuration)

	if err := repositories.CreateSession(tx, creds.ID, sessionToken, expirationTime); err != nil {
		log.Println("Error saving session to DB:", err)
		utils.HandleError(w, http.StatusInternalServerError, "Error creating session, try again later")
		return
	}

	if err := tx.Commit(); err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Server error")
		return
	}

	// set session cookie
	newCookie := http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  expirationTime,
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &newCookie)
}
