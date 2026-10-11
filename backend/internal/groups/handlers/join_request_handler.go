package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"zone/internal/auth"
	"zone/internal/database"
	"zone/internal/groups/repositories"
	"zone/internal/groups/services"
	notificationservice "zone/internal/notifications/services"
	"zone/internal/utils"
)

func Join_Request_handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.HandleError(w, 405, http.StatusText(405))
		return
	}
	user_id, err := auth.GetUserID(r)
	if err != nil {
		utils.HandleError(w, 401, http.StatusText(401))
		return
	}

	parts := strings.Split(r.URL.Path, "/")

	if len(parts) < 4 {
		utils.HandleError(w, 400, http.StatusText(400))
		return
	}

	group_id, err := strconv.Atoi(parts[2])
	if err != nil {
		utils.HandleError(w, 400, http.StatusText(400))
		return
	}

	if !services.JoinRequestValidator(group_id, user_id) {
		utils.HandleError(w, 400, http.StatusText(400))
		return
	}

	// Get the admin from the group_id

	creator_id, err := repositories.GetGroupCreator(group_id)
	if err != nil {
		utils.HandleError(w, 500, http.StatusText(500))
		return
	}

	if user_id == creator_id {
		utils.HandleError(w, 400, http.StatusText(400))
		return
	}

	// Use a transaction so the join request and its notification
	// are created atomically (all succeed or none is committed).
	tx, err := database.Database.Begin()
	if err != nil {
		utils.HandleError(w, 500, http.StatusText(500))
		return
	}
	defer tx.Rollback()

	request_id, err := repositories.SaveJoinRequest(tx, group_id, user_id)
	if err != nil {
		if err.Error() == "join request already exists" {
			utils.HandleError(w, http.StatusConflict, "join request already exists")
			return
		}
		utils.HandleError(w, 500, http.StatusText(500))
		return
	}

	fullname := repositories.GetUserFullName(user_id)

	msg := "New join request for your group " + fullname

	if err := tx.Commit(); err != nil {
		utils.HandleError(w, 500, http.StatusText(500))
		return
	}

	// creat the notifcation
	_, err = notificationservice.Notifier.Notify(creator_id, "join_request", &request_id, &user_id, msg)
	if err != nil {
		utils.HandleError(w, 500, http.StatusText(500))
		return
	}

	utils.RespondJSON(
		w,
		http.StatusCreated,
		map[string]string{
			"message": "join request  created successfully!",
		},
	)
}

func CancelJoinRequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.HandleError(w, 405, http.StatusText(405))
		return
	}
	user_id, err := auth.GetUserID(r)
	if err != nil {
		utils.HandleError(w, 401, http.StatusText(401))
		return
	}

	group_id, err := strconv.Atoi(r.PathValue("group_id"))
	if err != nil || group_id <= 0 {
		utils.HandleError(w, 400, "Invalid group id")
		return
	}

	request_id, err := repositories.GetJoinRequestIDByGroupAndUser(group_id, user_id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.HandleError(w, 404, "No pending join request")
			return
		}
		utils.HandleError(w, 500, http.StatusText(500))
		return
	}

	creator_id, err := repositories.GetGroupCreator(group_id)
	if err != nil {
		utils.HandleError(w, 500, http.StatusText(500))
		return
	}

	tx, err := database.Database.Begin()
	if err != nil {
		utils.HandleError(w, 500, http.StatusText(500))
		return
	}
	defer tx.Rollback()

	res, err := tx.Exec(
		`DELETE FROM group_requests
		WHERE id = ? AND user_id = ? AND group_id = ? AND status = 'pending'`,
		request_id, user_id, group_id,
	)
	if err != nil {
		utils.HandleError(w, 500, http.StatusText(500))
		return
	}
	affected, err := res.RowsAffected()
	if err != nil {
		utils.HandleError(w, 500, http.StatusText(500))
		return
	}
	if affected == 0 {
		utils.HandleError(w, 404, "No pending join request")
		return
	}

	if err := repositories.DeleteRequestNotification(tx, creator_id, request_id); err != nil {
		utils.HandleError(w, 500, http.StatusText(500))
		return
	}

	if err := tx.Commit(); err != nil {
		utils.HandleError(w, 500, http.StatusText(500))
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "join request cancelled successfully!"})
}
