package profile

import (
	"net/http"
	"strconv"

	"zone/internal/auth"
	"zone/internal/database"
	"zone/internal/models"
	postmodels "zone/internal/posts/models"
	postrepositories "zone/internal/posts/repositories"
	"zone/internal/utils"
)

func GetProfileAPI(w http.ResponseWriter, r *http.Request) {

	profileIDStr := r.PathValue("id")
	profileID, err := strconv.Atoi(profileIDStr)
	if err != nil || profileID <= 0 {
		utils.HandleError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	requesterID, _ := auth.GetUserID(r)

	var user models.User
	err = database.Database.QueryRow(`SELECT id, email, first_name, last_name, COALESCE(date_of_birth, ''), COALESCE(nickname, ''), COALESCE(avatar, ''), COALESCE(about_me, ''), is_private, created_at FROM users WHERE id = ?`, profileID).Scan(
		&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.DateOfBirth, &user.Nickname, &user.Avatar, &user.AboutMe, &user.IsPrivate, &user.CreatedAt,
	)
	if err != nil {
		utils.HandleError(w, http.StatusNotFound, "User not found")
		return
	}

	relation := "none"
	if requesterID == user.ID {
		relation = "self"
	} else if requesterID != 0 {
		var status string
		err := database.Database.QueryRow(`SELECT status FROM followers WHERE follower_id = ? AND followee_id = ?`, requesterID, user.ID).Scan(&status)
		if err == nil {
			relation = status
		}
	}

	if relation != "self" {
		user.Email = ""
	}

	var userPosts []postmodels.Post
	if user.IsPrivate && relation != "self" && relation != "accepted" {
		user.DateOfBirth = ""
		user.AboutMe = ""
		userPosts = []postmodels.Post{}
	} else {
		var err error
		userPosts, err = postrepositories.GetPostsByUser(profileID, requesterID)
		if err != nil || userPosts == nil {
			userPosts = []postmodels.Post{}
		}
	}

	var followersCount, followingCount int
	if err := database.Database.QueryRow(`
		SELECT
			COALESCE(SUM(followee_id = ?), 0),
			COALESCE(SUM(follower_id = ?), 0)
		FROM followers
		WHERE status = 'accepted' AND (followee_id = ? OR follower_id = ?)`,
		user.ID, user.ID, user.ID, user.ID,
	).Scan(&followersCount, &followingCount); err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Database error")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"user":            user,
		"relation":        relation,
		"followers_count": followersCount,
		"following_count": followingCount,
		"posts":           userPosts,
	})
}
