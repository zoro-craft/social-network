package profile

import (
	"net/http"
	"strconv"

	"zone/internal/auth"
	followersrepo "zone/internal/followers/repositories"
	postmodels "zone/internal/posts/models"
	postrepositories "zone/internal/posts/repositories"
	profilerepo "zone/internal/profile/repositories"
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

	user, err := profilerepo.GetUserByID(profileID)
	if err != nil {
		utils.HandleError(w, http.StatusNotFound, "User not found")
		return
	}

	relation := "none"
	if requesterID == user.ID {
		relation = "self"
	} else if requesterID != 0 {
		status, err := followersrepo.GetFollowStatus(requesterID, user.ID)
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

	followersCount, followingCount, err := followersrepo.GetFollowCounts(user.ID)
	if err != nil {
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
