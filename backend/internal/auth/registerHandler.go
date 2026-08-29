package auth

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"social-network/pkg/db"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{3,25}$`)
	emailRegex    = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)
)

type RegisterReq struct {
	Username    string `json:"username"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	DateOfBirth string `json:"date_of_birth"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(apiResponse{Status: "error", Message: "Method not allowed"})
		return
	}

	var req RegisterReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(apiResponse{Status: "error", Message: "Invalid request formatting"})
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	req.FirstName = strings.TrimSpace(req.FirstName)
	req.LastName = strings.TrimSpace(req.LastName)
	req.Email = strings.TrimSpace(req.Email)
	req.DateOfBirth = strings.TrimSpace(req.DateOfBirth)

	if !usernameRegex.MatchString(req.Username) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(apiResponse{Status: "error", Message: "Username must contain only letters, numbers, and underscores (3-25 chars)."})
		return
	}

	if req.DateOfBirth == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(apiResponse{Status: "error", Message: "Date of birth is required."})
		return
	}

	dob, err := time.Parse("2006-01-02", req.DateOfBirth)
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

	if req.FirstName == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(apiResponse{Status: "error", Message: "First name is required."})
		return
	}

	if req.LastName == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(apiResponse{Status: "error", Message: "Last name is required."})
		return
	}

	if !emailRegex.MatchString(req.Email) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(apiResponse{Status: "error", Message: "Please enter a valid email address."})
		return
	}

	var existingID int
	err = db.DB.QueryRow("SELECT id FROM users WHERE email = ? OR username = ?", req.Email, req.Username).Scan(&existingID)
	if err != sql.ErrNoRows {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(apiResponse{Status: "error", Message: "Email or username already exists"})
		return
	}

	if strings.TrimSpace(req.Password) == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(apiResponse{Status: "error", Message: "Password cannot be empty or consist only of spaces."})
		return
	}

	if req.Password != strings.TrimSpace(req.Password) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(apiResponse{Status: "error", Message: "Password cannot start or end with a space."})
		return
	}

	if len(req.Password) > 72 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(apiResponse{Status: "error", Message: "Password is too long! (72 characters max)"})
		return
	}
	if len(req.Password) < 8 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(apiResponse{Status: "error", Message: "Password is too short! (8 characters minimum)"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Bcrypt hash error:", err)
		json.NewEncoder(w).Encode(apiResponse{Status: "error", Message: "Error creating account, try again later"})
		return
	}

	_, err = db.DB.Exec(
		"INSERT INTO users (username, email, password, date_of_birth, first_name, last_name) VALUES (?, ?, ?, ?, ?, ?)",
		req.Username, req.Email, string(hashedPassword), dob.Format("2006-01-02"), req.FirstName, req.LastName,
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
