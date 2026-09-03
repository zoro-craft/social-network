package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"zone/internal/auth"
	"zone/internal/groups/repositories"
	"zone/internal/utils"
)
func Accept_Join_Request(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.HandleError(w, 405, http.StatusText(405))
		return
	}

	adminID, err := auth.GetUserIDFromSession(r)
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

	if joinRequest.Status != "pending" {
		utils.HandleError(w, 400, "request is not pending")
		return
	}

	creatorID, err := repositories.GetGroupCreator(joinRequest.GroupID)
	if err != nil {
		utils.HandleError(w, 500, http.StatusText(500))
		return
	}

	if creatorID != adminID {
		utils.HandleError(w, 403, http.StatusText(403))
		return
	}

	err = repositories.AddGroupMember(
		joinRequest.GroupID,
		joinRequest.UserID,
	)
	if err != nil {
		utils.HandleError(w, 500, http.StatusText(500))
		return
	}

	err = repositories.UpdateJoinRequestStatus(
		requestID,
		"accepted",
	)
	if err != nil {
		utils.HandleError(w, 500, http.StatusText(500))
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{
		"message": "request accepted successfully!",
	})
}

func Reject_Join_Request(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.HandleError(w, 405, http.StatusText(405))
		return
	}

	adminID, err := auth.GetUserIDFromSession(r)
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

	if joinRequest.Status != "pending" {
		utils.HandleError(w, 400, "request is not pending")
		return
	}

	creatorID, err := repositories.GetGroupCreator(joinRequest.GroupID)
	if err != nil {
		utils.HandleError(w, 500, http.StatusText(500))
		return
	}

	if creatorID != adminID {
		utils.HandleError(w, 403, http.StatusText(403))
		return
	}

	err = repositories.UpdateJoinRequestStatus(
		requestID,
		"rejected",
	)
	if err != nil {
		utils.HandleError(w, 500, http.StatusText(500))
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{
		"message": "request rejected successfully!",
	})
}
