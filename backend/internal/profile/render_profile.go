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
	var nickname, aboutMe, avatar string

	err = database.Database.QueryRow(`
		SELECT id, first_name, last_name, email, date_of_birth,
		       COALESCE(nickname, ''), COALESCE(about_me, ''), COALESCE(avatar, ''), is_public
		FROM users WHERE id = ?
	`, userID).Scan(&p.ID, &p.FirstName, &p.LastName, &p.Email, &p.DateOfBirth, &nickname, &aboutMe, &avatar, &p.IsPublic)
	if err != nil {
		utils.HandleError(w, http.StatusNotFound, "User not found")
		return
	}

	p.Nickname = nickname
	p.AboutMe = aboutMe
	p.Avatar = avatar

	database.Database.QueryRow(
		"SELECT COUNT(*) FROM followers WHERE followee_id = ? AND status = 'accepted'", userID,
	).Scan(&p.FollowersCount)

	database.Database.QueryRow(
		"SELECT COUNT(*) FROM followers WHERE follower_id = ? AND status = 'accepted'", userID,
	).Scan(&p.FollowingCount)

	isOwner := viewerID == p.ID
	canView := isOwner || p.IsPublic || isFollowing(viewerID, p.ID)

	if !canView {
		p.Email = ""
		p.DateOfBirth = ""
		p.AboutMe = ""
		p.Restricted = true
		p.Posts = []posts.Post{}
		utils.RespondJSON(w, http.StatusOK, p)
		return
	}

	p.Posts, err = posts.GetPostsByUser(userID, viewerID)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to load posts")
		return
	}

	utils.RespondJSON(w, http.StatusOK, p)
}

func isFollowing(viewerID, profileID int) bool {
	var exists bool
	database.Database.QueryRow(
		"SELECT EXISTS(SELECT * FROM followers WHERE follower_id = ? AND followee_id = ? AND status = 'accepted')",
		viewerID, profileID,
	).Scan(&exists)
	return exists
}