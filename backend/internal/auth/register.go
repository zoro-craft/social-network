package auth

import (
	"net/http"
	"regexp"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"zone/internal/database"
	"zone/internal/utils"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	user := User{
		FirstName:   strings.TrimSpace(r.FormValue("first_name")),
		LastName:    strings.TrimSpace(r.FormValue("last_name")),
		Email:       strings.TrimSpace(r.FormValue("email")),
		Password:    r.FormValue("password"),
		DateOfBirth: strings.TrimSpace(r.FormValue("date_of_birth")),
		Nickname:    strings.TrimSpace(r.FormValue("nickname")),
		AboutMe:     strings.TrimSpace(r.FormValue("about_me")),
		Avatar:      strings.TrimSpace(r.FormValue("avatar")),
	}
	confirmPassword := r.FormValue("confirm_password")

	// 1. Required fields
	if user.FirstName == "" || user.LastName == "" || user.Email == "" ||
		user.Password == "" || confirmPassword == "" || user.DateOfBirth == "" {
		utils.HandleError(w, http.StatusBadRequest, "First name, last name, email, password and date of birth are required")
		return
	}

	// 2. Name length
	if len(user.FirstName) < 2 || len(user.FirstName) > 50 ||
		len(user.LastName) < 2 || len(user.LastName) > 50 {
		utils.HandleError(w, http.StatusBadRequest, "Names must be between 2 and 50 characters")
		return
	}

	// 3. Email format
	if !emailRegex.MatchString(user.Email) || len(user.Email) > 100 {
		utils.HandleError(w, http.StatusBadRequest, "Invalid email address")
		return
	}

	// 4. Password strength + match
	if len(user.Password) < 6 || len(user.Password) > 72 {
		utils.HandleError(w, http.StatusBadRequest, "Password must be between 6 and 72 characters")
		return
	}
	if user.Password != confirmPassword {
		utils.HandleError(w, http.StatusBadRequest, "Password and confirm password do not match")
		return
	}

	// 5. Date of birth
	dob, err := time.Parse("2006-01-02", user.DateOfBirth)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Invalid date of birth format, expected YYYY-MM-DD")
		return
	}
	if dob.After(time.Now()) {
		utils.HandleError(w, http.StatusBadRequest, "Date of birth cannot be in the future")
		return
	}
	age := time.Since(dob).Hours() / 24 / 365.25
	if age < 13 {
		utils.HandleError(w, http.StatusBadRequest, "You must be at least 13 years old")
		return
	}

	// 6. Optional nickname length
	if user.Nickname != "" && (len(user.Nickname) < 2 || len(user.Nickname) > 50) {
		utils.HandleError(w, http.StatusBadRequest, "Nickname must be between 2 and 50 characters")
		return
	}

	// 7. Email uniqueness
	var emailExists bool
	if err := database.Database.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)", user.Email,
	).Scan(&emailExists); err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Database error checking email")
		return
	}
	if emailExists {
		utils.HandleError(w, http.StatusBadRequest, "Email already registered")
		return
	}

	// 8. Nickname uniqueness (only if provided)
	if user.Nickname != "" {
		var nicknameExists bool
		if err := database.Database.QueryRow(
			"SELECT EXISTS(SELECT 1 FROM users WHERE nickname = ?)", user.Nickname,
		).Scan(&nicknameExists); err != nil {
			utils.HandleError(w, http.StatusInternalServerError, "Database error checking nickname")
			return
		}
		if nicknameExists {
			utils.HandleError(w, http.StatusBadRequest, "Nickname already taken")
			return
		}
	}

	// 9. Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Password hashing error")
		return
	}

	// 10. Insert user
	res, err := database.Database.Exec(
		`INSERT INTO users (first_name, last_name, email, password, date_of_birth, nickname, about_me, avatar)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		user.FirstName, user.LastName, user.Email, string(hashedPassword),
		dob.Format("2006-01-02"),
		nullableString(user.Nickname), nullableString(user.AboutMe), nullableString(user.Avatar),
	)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Could not create account")
		return
	}

	id, _ := res.LastInsertId()

	utils.RespondJSON(w, http.StatusCreated, map[string]any{
		"message": "User created successfully!",
		"user_id": id,
	})
}

func nullableString(s string) any {
	if s == "" {
		return nil
	}
	return s
}
