package followers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"

	"zone/internal/auth"
	"zone/internal/database"
	followersrepo "zone/internal/followers/repositories"
	"zone/internal/groups/repositories"
	notificationservice "zone/internal/notifications/services"
	"zone/internal/utils"
)

var errUserNotFound = errors.New("follow target does not exist")

func ToggleFollowAPI(w http.ResponseWriter, r *http.Request) {
	followerID, err := auth.GetUserID(r)
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, "You must be logged in")
		return
	}

	followeeID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || followeeID <= 0 {
		utils.HandleError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	if followerID == followeeID {
		utils.HandleError(w, http.StatusBadRequest, "Cannot follow yourself")
		return
	}

	newStatus, err := toggleFollow(followerID, followeeID)
	if err != nil {
		if errors.Is(err, errUserNotFound) {
			utils.HandleError(w, http.StatusNotFound, "User not found")
			return
		}
		log.Println("toggle follow:", err)
		utils.HandleError(w, http.StatusInternalServerError, "Failed to update follow status")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"status": newStatus})
}

func toggleFollow(followerID, followeeID int) (string, error) {
	tx, err := database.Database.Begin()
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	isPrivate, hasStatus, prevStatus, err := followersrepo.GetFollowTargetState(tx, followerID, followeeID)
	if errors.Is(err, sql.ErrNoRows) {
		return "", errUserNotFound
	}
	if err != nil {
		return "", err
	}

	var newStatus string
	switch {
	case !hasStatus:
		newStatus = "accepted"
		if isPrivate {
			newStatus = "pending"
		}
		err = followersrepo.CreateFollow(tx, followerID, followeeID, newStatus)

	case prevStatus == "accepted", isPrivate:
		newStatus = "none"
		err = followersrepo.DeleteFollow(tx, followerID, followeeID)

	default:
		newStatus = "accepted"
		err = followersrepo.AcceptPendingFollow(tx, followerID, followeeID)
	}
	if err != nil {
		return "", err
	}

	if err := tx.Commit(); err != nil {
		return "", err
	}

	if newStatus == "pending" {
		message := "New follow request from " + repositories.GetUserFullName(followerID)
		_, _ = notificationservice.Notifier.NotifyFollowRequest(followeeID, followerID, message)
	}

	if newStatus == "none" && prevStatus == "pending" {
		_ = notificationservice.Notifier.RemoveFollowRequestNotification(followeeID, followerID)
	}

	return newStatus, nil
}
