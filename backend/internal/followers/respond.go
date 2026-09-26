package followers

import (
	"net/http"
	"strconv"

	"zone/internal/auth"
	"zone/internal/database"
	followersrepo "zone/internal/followers/repositories"
	notificationservice "zone/internal/notifications/services"
	"zone/internal/utils"
)

func RespondRequestAPI(w http.ResponseWriter, r *http.Request) {
	followerID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || followerID <= 0 {
		utils.HandleError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	followeeID, err := auth.GetUserID(r)
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	action := r.PathValue("action")
	if action != "accept" && action != "decline" {
		utils.HandleError(w, http.StatusBadRequest, "Invalid action")
		return
	}

	tx, err := database.Database.Begin()
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to process request")
		return
	}
	defer tx.Rollback()

	var rowsAffected int64
	if action == "accept" {
		rowsAffected, err = followersrepo.AcceptFollowRequest(tx, followerID, followeeID)
	} else {
		rowsAffected, err = followersrepo.DeclineFollowRequest(tx, followerID, followeeID)
	}

	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to update request")
		return
	}

	if rowsAffected == 0 {
		utils.HandleError(w, http.StatusNotFound, "Request not found")
		return
	}

	if err := tx.Commit(); err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to process request")
		return
	}

	_ = notificationservice.Notifier.MarkReadOnAction(followeeID, "follow_request", followerID)

	utils.RespondJSON(w, http.StatusOK, map[string]string{"status": "success"})
}
