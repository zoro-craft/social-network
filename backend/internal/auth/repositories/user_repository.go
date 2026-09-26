package repositories

import (
	"database/sql"

	"zone/internal/database"
)

type UserInsert struct {
	Email          string
	HashedPassword string
	DateOfBirth    string
	FirstName      string
	LastName       string
	Nickname       any
	AboutMe        string
	AvatarPath     string
}

type Credentials struct {
	ID       int
	Password string
	Nickname string
	Email    string
}

type Profile struct {
	Nickname  string
	Email     string
	FirstName string
	LastName  string
	Avatar    string
}

// checks email existence
func EmailExists(tx *sql.Tx, email string) (bool, error) {
	var existingID int

	err := tx.QueryRow("SELECT id FROM users WHERE email = ?", email).Scan(&existingID)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

// checks nickname availability
func NicknameTaken(nickname string, excludeID int) (bool, error) {
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

// inserts a new user
func CreateUser(tx *sql.Tx, u UserInsert) error {
	_, err := tx.Exec(
		"INSERT INTO users (email, password, date_of_birth, first_name, last_name, nickname, about_me, avatar) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		u.Email, u.HashedPassword, u.DateOfBirth, u.FirstName, u.LastName, u.Nickname, u.AboutMe, u.AvatarPath,
	)

	return err
}

// loads login credentials
func CredentialsByIdentifier(tx *sql.Tx, identify string) (Credentials, error) {
	var creds Credentials

	err := tx.QueryRow(`
        SELECT id, password, COALESCE(nickname, ''), email FROM users 
        WHERE email = ? OR (nickname = ? AND nickname != '')
    `, identify, identify).Scan(&creds.ID, &creds.Password, &creds.Nickname, &creds.Email)
	if err != nil {
		return creds, err
	}

	return creds, nil
}

// loads user profile fields
func ProfileByID(userID int) (Profile, error) {
	var p Profile

	err := database.Database.QueryRow(
		"SELECT COALESCE(nickname, ''), email, COALESCE(first_name, ''), COALESCE(last_name, ''), COALESCE(avatar, '') FROM users WHERE id = ?",
		userID,
	).Scan(&p.Nickname, &p.Email, &p.FirstName, &p.LastName, &p.Avatar)
	if err != nil {
		return p, err
	}

	return p, nil
}
