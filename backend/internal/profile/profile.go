package profile

import (
	"encoding/json"
	"net/http"
	"strconv"

	"zone/internal/auth"
	"zone/internal/database"
	"zone/internal/models"
)

func GetProfileAPI(w http.ResponseWriter, r *http.Request) {

	profileIDStr := r.PathValue("id")
	profileID, _ := strconv.Atoi(profileIDStr)

	requesterID, _ := auth.GetUserID(r)

	var user models.User
	err := database.Database.QueryRow(`SELECT id, email, first_name, last_name, nickname, avatar, about_me, is_private, created_at FROM users WHERE id = ?`, profileID).Scan(
		&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.Nickname, &user.Avatar, &user.AboutMe, &user.IsPrivate, &user.CreatedAt,
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
	if user.IsPrivate && relation != "self" && relation != "accepted" {
		user.Email = ""
		user.AboutMe = ""
	}

	var followersCount, followingCount int
	database.Database.QueryRow(`SELECT count(*) FROM followers WHERE followee_id = ? AND status = 'accepted'`, user.ID).Scan(&followersCount)
	database.Database.QueryRow(`SELECT count(*) FROM followers WHERE follower_id = ? AND status = 'accepted'`, user.ID).Scan(&followingCount)
	response := models.UserProfile{
		User:      user,
		Followers: followersCount,
		Following: followingCount,
		Relation:  relation,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
