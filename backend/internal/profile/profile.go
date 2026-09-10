package profile

import (
	"net/http"
	"strconv"

	"zone/internal/auth"
	"zone/internal/database"
	"zone/internal/models"
	"zone/internal/posts"
	"zone/internal/utils"
)

func GetProfileAPI(w http.ResponseWriter, r *http.Request) {

	profileIDStr := r.PathValue("id")
	profileID, _ := strconv.Atoi(profileIDStr)

	requesterID, _ := auth.GetUserID(r)

	var user models.User
	err := database.Database.QueryRow(`SELECT id, email, first_name, last_name, COALESCE(date_of_birth, ''), nickname, avatar, about_me, is_private, created_at FROM users WHERE id = ?`, profileID).Scan(
		&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.DateOfBirth, &user.Nickname, &user.Avatar, &user.AboutMe, &user.IsPrivate, &user.CreatedAt,
	)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
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

	var userPosts []posts.Post
	if user.IsPrivate && relation != "self" && relation != "accepted" {
		user.Email = ""
		user.DateOfBirth = ""
		user.AboutMe = ""
		userPosts = []posts.Post{}
	} else {
		var err error
		userPosts, err = posts.GetPostsByUser(profileID, requesterID)
		if err != nil || userPosts == nil {
			userPosts = []posts.Post{}
		}
	}

	var followersCount, followingCount int
	database.Database.QueryRow(`SELECT count(*) FROM followers WHERE followee_id = ? AND status = 'accepted'`, user.ID).Scan(&followersCount)
	database.Database.QueryRow(`SELECT count(*) FROM followers WHERE follower_id = ? AND status = 'accepted'`, user.ID).Scan(&followingCount)

	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"user":            user,
		"relation":        relation,
		"followers_count": followersCount,
		"following_count": followingCount,
		"posts":           userPosts,
	})
}
