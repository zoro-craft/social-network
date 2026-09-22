package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"zone/internal/auth"
	"zone/internal/database"
	"zone/internal/groups/repositories"

	notificationservice "zone/internal/notifications/services"
	"zone/internal/utils"
	ws "zone/internal/websocket"
)

func Accept_Join_Request(w http.ResponseWriter, r *http.Request) {
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

	requestID, err := strconv.Atoi(parts[2])
	if err != nil {
		utils.HandleError(w, 400, http.StatusText(400))
		return
	}

	joinRequest, err := repositories.GetJoinRequest(requestID)
	if err != nil {
		utils.HandleError(w, 404, http.StatusText(404))
		return
	}

	creatorID, err := repositories.GetGroupCreator(joinRequest.GroupID)
	if err != nil {
		utils.HandleError(w, 500, http.StatusText(500))
		return
	}

	if creatorID != user_id {
		utils.HandleError(w, 403, http.StatusText(403))
		return
	}

	// Use a transaction so adding the member and deleting the request are atomic.
	tx, err := database.Database.Begin()
	if err != nil {
		utils.HandleError(w, 500, http.StatusText(500))
		return
	}
	defer tx.Rollback()

	if repositories.IsMember(joinRequest.UserID, joinRequest.GroupID) {
		repositories.DeleteRequest(tx, requestID)
		repositories.DeleteRequestNotification(tx, user_id, requestID)
		tx.Commit()
		return
	}

	err = repositories.AddGroupMember(tx, joinRequest.GroupID, joinRequest.UserID)
	if err != nil {
		utils.HandleError(w, 500, http.StatusText(500))
		return
	}

	err = repositories.DeleteRequest(tx, requestID)
	if err != nil {
		utils.HandleError(w, 500, http.StatusText(500))
		return
	}

	if err := tx.Commit(); err != nil {
		utils.HandleError(w, 500, http.StatusText(500))
		return
	}

	// Resolving the join request marks its notification read.
	_ = notificationservice.Notifier.MarkReadOnAction(user_id, "join_request", requestID)

	ws.DefaultHub.Subscribe(joinRequest.UserID, joinRequest.GroupID) // subscribe the new member to the group room

	utils.RespondJSON(w, http.StatusOK, map[string]string{
		"message": "request accepted successfully!",
	})
}

func Reject_Join_Request(w http.ResponseWriter, r *http.Request) {
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

	requestID, err := strconv.Atoi(parts[2])
	if err != nil {
		utils.HandleError(w, 400, http.StatusText(400))
		return
	}

	joinRequest, err := repositories.GetJoinRequest(requestID)
	if err != nil {
		utils.HandleError(w, 404, http.StatusText(404))
		return
	}

	creatorID, err := repositories.GetGroupCreator(joinRequest.GroupID)
	if err != nil {
		utils.HandleError(w, 500, http.StatusText(500))
		return
	}

	if creatorID != user_id {
		utils.HandleError(w, 403, http.StatusText(403))
		return
	}

	// Use a transaction so deleting the request and marking it read are atomic.
	tx, err := database.Database.Begin()
	if err != nil {
		utils.HandleError(w, 500, http.StatusText(500))
		return
	}
	defer tx.Rollback()

	err = repositories.DeleteRequest(tx, requestID)
	if err != nil {
		utils.HandleError(w, 500, http.StatusText(500))
		return
	}

	if err := tx.Commit(); err != nil {
		utils.HandleError(w, 500, http.StatusText(500))
		return
	}

	// Resolving the join request marks its notification read.
	_ = notificationservice.Notifier.MarkReadOnAction(user_id, "join_request", requestID)

	utils.RespondJSON(w, http.StatusOK, map[string]string{
		"message": "request rejected successfully!",
	})
}