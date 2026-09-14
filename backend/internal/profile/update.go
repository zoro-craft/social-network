package profile

import (
	"net/http"
	"strings"

	"zone/internal/auth"
	"zone/internal/database"
	"zone/internal/utils"
)

func SwitchStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID, err := auth.GetUserID(r)
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	_, err = database.Database.Exec(
		"UPDATE users SET is_private = NOT is_private WHERE id = ?", userID,
	)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to update profile")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Profile visibility updated"})
}

func EditProfileAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID, err := auth.GetUserID(r)
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 20<<20)
	err = r.ParseMultipartForm(20 << 20)
	if err != nil {
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
	password := r.FormValue("password")

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

	if password != "" {
		if err := auth.ValidatePassword(password); err != nil {
			utils.HandleError(w, http.StatusBadRequest, err.Error())
			return
		}
	}

	var nicknameValue any
	if nickname != "" {
		taken, err := auth.ValidateNickname(nickname, userID)
		if taken {
			utils.HandleError(w, http.StatusConflict, "Nickname already exists")
			return
		} else if err != nil {
			utils.HandleError(w, http.StatusInternalServerError, "Failed to update profile")
			return
		}
		nicknameValue = nickname
	}

	removeAvatar := r.FormValue("remove_avatar") == "on"

	// avatar handling
	var avatarPath string
	updateAvatar := false

	file, header, err := r.FormFile("avatar")
	if err == nil {
		defer file.Close()
		filename, saveErr := utils.SaveUploadFile(file, header.Filename, "./uploads/avatars")
		if saveErr == nil {
			avatarPath = "/uploads/avatars/" + filename
			updateAvatar = true
		}
	} else if removeAvatar {
		avatarPath = ""
		updateAvatar = true
	}

	if updateAvatar {
		_, err = database.Database.Exec(`
			UPDATE users 
			SET first_name = ?, last_name = ?, email = ?, date_of_birth = ?, nickname = ?, about_me = ?, is_private = ?, avatar = ? 
			WHERE id = ?`,
			firstName, lastName, email, dob, nicknameValue, aboutMe, isPrivate, avatarPath, userID,
		)
	} else {
		_, err = database.Database.Exec(`
			UPDATE users 
			SET first_name = ?, last_name = ?, email = ?, date_of_birth = ?, nickname = ?, about_me = ?, is_private = ? 
			WHERE id = ?`,
			firstName, lastName, email, dob, nicknameValue, aboutMe, isPrivate, userID,
		)
	}

	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to update profile")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Profile updated successfully"})
}
