package handlers

import (
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
	// creat the notifcation
	_, err = notificationservice.Notifier.Notify(creator_id, "join_request", &request_id, &user_id, msg)
	if err != nil {
		utils.HandleError(w, 500, http.StatusText(500))
		return
	}

	if err := tx.Commit(); err != nil {
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
