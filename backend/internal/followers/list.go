package followers

import (
	"net/http"
	"strconv"

	"zone/internal/auth"
	"zone/internal/followers/models"
	followersrepo "zone/internal/followers/repositories"
	"zone/internal/utils"
)

type UserItem = models.UserItem

func canViewList(viewerID, profileID int) bool {
	if viewerID == profileID {
		return true
	}

	isPrivate, err := followersrepo.IsPrivateProfile(profileID)
	if err != nil {
		return false
	}

	if !isPrivate {
		return true
	}

	isFollower, err := followersrepo.IsAcceptedFollower(viewerID, profileID)
	return err == nil && isFollower
}

func GetFollowersAPI(w http.ResponseWriter, r *http.Request) {
	listHandler(w, r, true)
}

func GetFollowingAPI(w http.ResponseWriter, r *http.Request) {
	listHandler(w, r, false)
}

func listHandler(w http.ResponseWriter, r *http.Request, listFollowers bool) {
	profileID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || profileID <= 0 {
		utils.HandleError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	viewerID, _ := auth.GetUserID(r)

	if !canViewList(viewerID, profileID) {
		utils.RespondJSON(w, http.StatusOK, []UserItem{})
		return
	}

	listUsers(w, profileID, listFollowers)
}

func listUsers(w http.ResponseWriter, profileID int, listFollowers bool) {
	users, err := followersrepo.ListConnections(profileID, listFollowers)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Database error")
		return
	}

	utils.RespondJSON(w, http.StatusOK, users)
}
