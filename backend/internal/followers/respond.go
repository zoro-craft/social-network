package followers

import (
	"database/sql"
	"net/http"
	"strconv"

	"zone/internal/auth"
	"zone/internal/database"
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

	var result sql.Result
	if action == "accept" {
		result, err = tx.Exec(
			"UPDATE followers SET status = 'accepted' WHERE follower_id = ? AND followee_id = ? AND status = 'pending'",
			followerID, followeeID,
		)
	} else {
		result, err = tx.Exec(
			"DELETE FROM followers WHERE follower_id = ? AND followee_id = ? AND status = 'pending'",
			followerID, followeeID,
		)
	}

	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to update request")
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
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
