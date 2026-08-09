package profile

import (
	"net/http"
	"strconv"

	"zone/internal/auth"
	"zone/internal/database"
	"zone/internal/posts"
	"zone/internal/utils"
)

func GetProfileAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	idStr := r.PathValue("id")
	userID, err := strconv.Atoi(idStr)
	if err != nil || userID <= 0 {
		utils.HandleError(w, http.StatusBadRequest, "Invalid user id")
		return
	}

	viewerID, _ := auth.GetUserIDFromSession(r)

	var p Profile
	var nickname, aboutMe, avatar *string

	err = database.Database.QueryRow(`
		SELECT id, first_name, last_name, email, date_of_birth, nickname, about_me, avatar, is_public
		FROM users WHERE id = ?
	`, userID).Scan(&p.ID, &p.FirstName, &p.LastName, &p.Email, &p.DateOfBirth, &nickname, &aboutMe, &avatar, &p.IsPublic)
	if err != nil {
		utils.HandleError(w, http.StatusNotFound, "User not found")
		return
	}

	if nickname != nil {
		p.Nickname = *nickname
	}
	if aboutMe != nil {
		p.AboutMe = *aboutMe
	}
	if avatar != nil {
		p.Avatar = *avatar
	}

	p.Posts, err = posts.GetPostsByUser(userID, viewerID)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to load posts")
		return
	}

	utils.RespondJSON(w, http.StatusOK, p)
}
