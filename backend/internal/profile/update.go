package profile

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strings"

	"zone/internal/auth"
	"zone/internal/database"
	"zone/internal/utils"
)

func EditProfileAPI(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserID(r)
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 5<<20)
	if err := r.ParseMultipartForm(utils.MaxUploadSize); err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Form too large or invalid")
		return
	}

	isPrivate := r.FormValue("is_private") == "on"
	firstName := strings.TrimSpace(r.FormValue("first_name"))
	lastName := strings.TrimSpace(r.FormValue("last_name"))
	email := r.FormValue("email")
	dob := strings.TrimSpace(r.FormValue("date_of_birth"))
	nickname := strings.TrimSpace(r.FormValue("nickname"))
	aboutMe := r.FormValue("about_me")

	if _, err := auth.ValidateDOB(dob); err != nil {
		utils.HandleError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := auth.ValidateName(firstName, lastName); err != nil {
		utils.HandleError(w, http.StatusBadRequest, err.Error())
		return
	}
	if !auth.ValidateEmail(email) {
		utils.HandleError(w, http.StatusBadRequest, "Please enter a valid email address.")
		return
	}

	var emailTaken bool
	if err := database.Database.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM users WHERE email = ? AND id != ?)", email, userID,
	).Scan(&emailTaken); err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Database error")
		return
	}
	if emailTaken {
		utils.HandleError(w, http.StatusConflict, "Email already exists")
		return
	}

	if nickname != "" {
		taken, err := auth.ValidateNickname(nickname, userID)
		if err != nil {
			utils.HandleError(w, http.StatusInternalServerError, "Failed to check nickname")
			return
		}
		if taken {
			utils.HandleError(w, http.StatusConflict, "Nickname already exists")
			return
		}
	}

	var avatarPath string
	updateAvatar := false
	avatarPath, err = utils.SaveImageField(r, "avatar", "avatars", "/uploads/avatars/")
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Only JPEG, PNG, and GIF images are allowed")
		return
	}
	if avatarPath != "" {
		updateAvatar = true
	} else if r.FormValue("remove_avatar") == "on" {
		avatarPath = ""
		updateAvatar = true
	}

	tx, err := database.Database.Begin()
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Database error")
		return
	}
	defer tx.Rollback()

	var wasPrivate bool
	if err := tx.QueryRow("SELECT is_private FROM users WHERE id = ?", userID).Scan(&wasPrivate); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.HandleError(w, http.StatusNotFound, "User not found")
			return
		}
		utils.HandleError(w, http.StatusInternalServerError, "Database error")
		return
	}

	if _, err := tx.Exec(`
		UPDATE users SET
			first_name = ?, last_name = ?, email = ?, date_of_birth = ?,
			nickname = CASE WHEN ? <> '' THEN ? ELSE nickname END,
			about_me = ?, is_private = ?,
			avatar = CASE WHEN ? THEN ? ELSE avatar END
		WHERE id = ?`,
		firstName, lastName, email, dob,
		nickname, nickname, aboutMe, isPrivate,
		updateAvatar, avatarPath, userID,
	); err != nil {
		log.Println("update profile:", err)
		utils.HandleError(w, http.StatusInternalServerError, "Failed to update profile")
		return
	}

	if wasPrivate && !isPrivate {
		if _, err := tx.Exec(
			"UPDATE followers SET status = 'accepted' WHERE followee_id = ? AND status = 'pending'",
			userID,
		); err != nil {
			log.Println("accept pending followers:", err)
			utils.HandleError(w, http.StatusInternalServerError, "Failed to update profile")
			return
		}
	}

	if err := tx.Commit(); err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to update profile")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Profile updated successfully"})
}
