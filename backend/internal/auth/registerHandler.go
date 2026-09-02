package auth

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"zone/internal/database"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	authEmailRegex = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(apiResponse{Status: "error", Message: "Method not allowed"})
		return
	}

	contentType := r.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "multipart/form-data") {
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(apiResponse{Status: "error", Message: "Invalid form data or file too large"})
			return
		}
	} else {
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(apiResponse{Status: "error", Message: "Invalid form data"})
			return
		}
	}

	firstName := strings.TrimSpace(r.FormValue("first_name"))
	lastName := strings.TrimSpace(r.FormValue("last_name"))
	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password")
	dateOfBirth := strings.TrimSpace(r.FormValue("date_of_birth"))
	nickname := strings.TrimSpace(r.FormValue("nickname"))
	aboutMe := strings.TrimSpace(r.FormValue("about_me"))

	if dateOfBirth == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(apiResponse{Status: "error", Message: "Date of birth is required."})
		return
	}

	dob, err := time.Parse("2006-01-02", dateOfBirth)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(apiResponse{Status: "error", Message: "Invalid date of birth format. Must be YYYY-MM-DD."})
		return
	}

	if now := time.Now(); dob.After(now.AddDate(-13, 0, 0)) || dob.Before(now.AddDate(-120, 0, 0)) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(apiResponse{Status: "error", Message: "You must be between 13 and 120 years old to register."})
		return
	}

	if firstName == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(apiResponse{Status: "error", Message: "First name is required."})
		return
	}

	if lastName == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(apiResponse{Status: "error", Message: "Last name is required."})
		return
	}

	if !authEmailRegex.MatchString(email) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(apiResponse{Status: "error", Message: "Please enter a valid email address."})
		return
	}

	var existingID int
	err = database.Database.QueryRow("SELECT id FROM users WHERE email = ?", email).Scan(&existingID)
	if err == nil {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(apiResponse{Status: "error", Message: "Email already exists"})
		return
	} else if err != sql.ErrNoRows {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Database query error checking existing email:", err)
		json.NewEncoder(w).Encode(apiResponse{Status: "error", Message: "Database error, please try again"})
		return
	}

	if nickname != "" {
		var existingNicknameID int
		err = database.Database.QueryRow("SELECT id FROM users WHERE nickname = ?", nickname).Scan(&existingNicknameID)
		if err == nil {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(apiResponse{Status: "error", Message: "Nickname already exists"})
			return
		} else if err != sql.ErrNoRows {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("Database query error checking existing nickname:", err)
			json.NewEncoder(w).Encode(apiResponse{Status: "error", Message: "Database error, please try again"})
			return
		}
	}

	if strings.TrimSpace(password) == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(apiResponse{Status: "error", Message: "Password cannot be empty or consist only of spaces."})
		return
	}

	if password != strings.TrimSpace(password) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(apiResponse{Status: "error", Message: "Password cannot start or end with a space."})
		return
	}

	if len(password) > 72 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(apiResponse{Status: "error", Message: "Password is too long! (72 characters max)"})
		return
	}
	if len(password) < 8 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(apiResponse{Status: "error", Message: "Password is too short! (8 characters minimum)"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Bcrypt hash error:", err)
		json.NewEncoder(w).Encode(apiResponse{Status: "error", Message: "Error creating account, try again later"})
		return
	}

	var avatarPath string
	if strings.HasPrefix(contentType, "multipart/form-data") {
		file, handler, err := r.FormFile("avatar")
		if err == nil {
			defer file.Close()
			os.MkdirAll("uploads/avatars", os.ModePerm)
			ext := filepath.Ext(handler.Filename)
			newFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
			filePath := filepath.Join("uploads/avatars", newFileName)

			dst, err := os.Create(filePath)
			if err == nil {
				defer dst.Close()
				io.Copy(dst, file)
				avatarPath = "/" + filePath
			}
		}
	} else {
		avatarPath = strings.TrimSpace(r.FormValue("avatar"))
	}

	_, err = database.Database.Exec(
		"INSERT INTO users (email, password, date_of_birth, first_name, last_name, nickname, about_me, avatar) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		email, string(hashedPassword), dob.Format("2006-01-02"), firstName, lastName, nickname, aboutMe, avatarPath,
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Database insert error:", err)
		json.NewEncoder(w).Encode(apiResponse{Status: "error", Message: "Could not create account"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(apiResponse{
		Status:  "success",
		Message: "Account created successfully. Please login.",
	})
}
