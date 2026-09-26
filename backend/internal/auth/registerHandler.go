package auth

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
	"zone/internal/utils"

	"zone/internal/auth/repositories"
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

	emailExists, err := repositories.EmailExists(tx, email)
	if err != nil {
		log.Println("Database query error checking existing email:", err)
		utils.HandleError(w, http.StatusInternalServerError, "Database error, please try again")
		return
	}
	if emailExists {
		utils.HandleError(w, http.StatusConflict, "Email already exists")
		return
	}

	var nicknameValue any
	if nickname != "" {
		if utf8.RuneCountInString(nickname) < 3 {
			utils.HandleError(w, http.StatusBadRequest, "Nickname must be at least 3 characters long.")
			return
		}
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

	err = repositories.CreateUser(tx, repositories.UserInsert{
		Email:          email,
		HashedPassword: string(hashedPassword),
		DateOfBirth:    dob.Format("2006-01-02"),
		FirstName:      firstName,
		LastName:       lastName,
		Nickname:       nicknameValue,
		AboutMe:        aboutMe,
		AvatarPath:     avatarPath,
	})
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
	if utf8.RuneCountInString(firstName) < 3 {
		return fmt.Errorf("First name must be at least 3 characters long.")
	}
	if utf8.RuneCountInString(lastName) < 3 {
		return fmt.Errorf("Last name must be at least 3 characters long.")
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

// checks nickname availability
func ValidateNickname(nickname string, excludeID int) (bool, error) {
	return repositories.NicknameTaken(nickname, excludeID)
}
