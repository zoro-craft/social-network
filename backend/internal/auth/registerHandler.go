package auth

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
	"zone/internal/utils"

	"zone/internal/database"

	"golang.org/x/crypto/bcrypt"
)

var (
	authEmailRegex = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 5<<20)
	contentType := r.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "multipart/form-data") {
		if err := r.ParseMultipartForm(utils.MaxUploadSize); err != nil {
			utils.HandleError(w, http.StatusBadRequest, "Invalid form data or file too large")
			return
		}
	} else {
		if err := r.ParseForm(); err != nil {
			utils.HandleError(w, http.StatusBadRequest, "Invalid form data")
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

	dob, err := ValidateDOB(dateOfBirth)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := ValidateName(firstName, lastName); err != nil {
		utils.HandleError(w, http.StatusBadRequest, err.Error())
		return
	}

	if !ValidateEmail(email) {
		utils.HandleError(w, http.StatusBadRequest, "Please enter a valid email address.")
		return
	}

	tx, err := database.Database.Begin()
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Database error, please try again")
		return
	}
	defer tx.Rollback()

	var existingID int
	err = tx.QueryRow("SELECT id FROM users WHERE email = ?", email).Scan(&existingID)
	if err == nil {
		utils.HandleError(w, http.StatusConflict, "Email already exists")
		return
	} else if err != sql.ErrNoRows {
		log.Println("Database query error checking existing email:", err)
		utils.HandleError(w, http.StatusInternalServerError, "Database error, please try again")
		return
	}

	var nicknameValue any
	if nickname != "" {
		taken, err := ValidateNickname(nickname, 0)
		if taken {
			utils.HandleError(w, http.StatusConflict, "Nickname already exists")
			return
		} else if err != nil {
			log.Println("Database query error checking existing nickname:", err)
			utils.HandleError(w, http.StatusInternalServerError, "Database error, please try again")
			return
		}
		nicknameValue = nickname
	}

	if err := ValidatePassword(password); err != nil {
		utils.HandleError(w, http.StatusBadRequest, err.Error())
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Bcrypt hash error:", err)
		utils.HandleError(w, http.StatusInternalServerError, "Error creating account, try again later")
		return
	}

	avatarPath, err := utils.SaveImageField(r, "avatar", "avatars", "/uploads/avatars/")
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Invalid image")
		return
	}

	err = createUser(tx, email, string(hashedPassword), dob.Format("2006-01-02"), firstName, lastName, nicknameValue, aboutMe, avatarPath)
	if err != nil {
		log.Println("Database insert error:", err)
		utils.HandleError(w, http.StatusInternalServerError, "Could not create account")
		return
	}

	if err := tx.Commit(); err != nil {
		log.Println("Transaction commit error:", err)
		utils.HandleError(w, http.StatusInternalServerError, "Could not create account")
		return
	}

	utils.RespondJSON(w, http.StatusCreated, map[string]string{"status": "success", "message": "Account created successfully. Please login."})
}

func ValidateEmail(email string) bool {
	return authEmailRegex.MatchString(email)
}

func ValidateDOB(dateOfBirth string) (time.Time, error) {
	if dateOfBirth == "" {
		return time.Time{}, fmt.Errorf("Date of birth is required.")
	}

	dob, err := time.Parse("2006-01-02", dateOfBirth)
	if err != nil {
		return time.Time{}, fmt.Errorf("Invalid date of birth format. Must be YYYY-MM-DD.")
	}

	if now := time.Now(); dob.After(now.AddDate(-13, 0, 0)) || dob.Before(now.AddDate(-120, 0, 0)) {
		return time.Time{}, fmt.Errorf("You must be between 13 and 120 years old to register.")
	}

	return dob, nil
}

func ValidateName(firstName, lastName string) error {
	if firstName == "" || lastName == "" {
		return fmt.Errorf("First & last name  is required.")
	}
	return nil
}

func ValidatePassword(password string) error {
	if strings.TrimSpace(password) == "" {
		return fmt.Errorf("Password cannot be empty or consist only of spaces.")
	}
	if password != strings.TrimSpace(password) {
		return fmt.Errorf("Password cannot start or end with a space.")
	}
	if len(password) > 72 {
		return fmt.Errorf("Password is too long! (72 characters max)")
	}
	if len(password) < 8 {
		return fmt.Errorf("Password is too short! (8 characters minimum)")
	}
	return nil
}

func ValidateNickname(nickname string, excludeID int) (bool, error) {
	var query string
	var args []any
	if excludeID > 0 {
		query = "SELECT id FROM users WHERE nickname = ? AND id != ?"
		args = []any{nickname, excludeID}
	} else {
		query = "SELECT id FROM users WHERE nickname = ?"
		args = []any{nickname}
	}

	var existingID int
	err := database.Database.QueryRow(query, args...).Scan(&existingID)
	if err == nil {
		return true, nil
	}
	if err != sql.ErrNoRows {
		return false, err
	}
	return false, nil
}

func createUser(tx *sql.Tx, email, hashedPassword, dob, firstName, lastName string, nickname any, aboutMe, avatarPath string) error {
	_, err := tx.Exec(
		"INSERT INTO users (email, password, date_of_birth, first_name, last_name, nickname, about_me, avatar) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		email, hashedPassword, dob, firstName, lastName, nickname, aboutMe, avatarPath,
	)
	return err
}
