package repositories

import (
	"database/sql"

	"zone/internal/database"
	"zone/internal/models"
)

type ProfileUpdate struct {
	UserID       int
	FirstName    string
	LastName     string
	Email        string
	DateOfBirth  string
	Nickname     string
	AboutMe      string
	IsPrivate    bool
	UpdateAvatar bool
	AvatarPath   string
}

func GetUserByID(profileID int) (models.User, error) {
	var user models.User

	err := database.Database.QueryRow(`SELECT id, email, first_name, last_name, COALESCE(date_of_birth, ''), COALESCE(nickname, ''), COALESCE(avatar, ''), COALESCE(about_me, ''), is_private, created_at FROM users WHERE id = ?`, profileID).Scan(
		&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.DateOfBirth, &user.Nickname, &user.Avatar, &user.AboutMe, &user.IsPrivate, &user.CreatedAt,
	)
	if err != nil {
		return user, err
	}

	return user, nil
}

func EmailTakenByOther(email string, userID int) (bool, error) {
	var emailTaken bool

	err := database.Database.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM users WHERE email = ? AND id != ?)", email, userID,
	).Scan(&emailTaken)
	if err != nil {
		return false, err
	}

	return emailTaken, nil
}

func GetPrivacyFlag(tx *sql.Tx, userID int) (bool, error) {
	var isPrivate bool

	err := tx.QueryRow("SELECT is_private FROM users WHERE id = ?", userID).Scan(&isPrivate)
	if err != nil {
		return false, err
	}

	return isPrivate, nil
}

func UpdateProfile(tx *sql.Tx, p ProfileUpdate) error {
	_, err := tx.Exec(`
		UPDATE users SET
			first_name = ?, last_name = ?, email = ?, date_of_birth = ?,
			nickname = CASE WHEN ? <> '' THEN ? ELSE nickname END,
			about_me = ?, is_private = ?,
			avatar = CASE WHEN ? THEN ? ELSE avatar END
		WHERE id = ?`,
		p.FirstName, p.LastName, p.Email, p.DateOfBirth,
		p.Nickname, p.Nickname, p.AboutMe, p.IsPrivate,
		p.UpdateAvatar, p.AvatarPath, p.UserID,
	)

	return err
}
