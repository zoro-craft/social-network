package handlers

import (
	"net/http"
	"strconv"

	"zone/internal/auth"
	"zone/internal/groups/repositories"
	"zone/internal/groups/services"
	"zone/internal/utils"
)

// RequestToJoinHandler lets a non-member request to join a group.
func RequestToJoinHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID, err := auth.GetUserIDFromSession(r)
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, "You must be logged in")
		return
	}

	groupID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || groupID <= 0 {
		utils.HandleError(w, http.StatusBadRequest, "Invalid group id")
		return
	}

	if !services.GroupExists(groupID) {
		utils.HandleError(w, http.StatusNotFound, "Group not found")
		return
	}

	if services.IsGroupMember(groupID, userID) || services.IsGroupCreator(groupID, userID) {
		utils.HandleError(w, http.StatusBadRequest, "You are already a member of this group")
		return
	}

	pending, err := repositories.HasPendingRequest(groupID, userID)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to check request")
		return
	}
	if pending {
		utils.HandleError(w, http.StatusBadRequest, "You already have a pending request")
		return
	}

	if err := repositories.CreateRequest(groupID, userID); err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to send join request")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Join request sent"})
}

// GroupRequestsHandler lists join requests. Only the creator can do this.
func GroupRequestsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID, err := auth.GetUserIDFromSession(r)
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, "You must be logged in")
		return
	}

	groupID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || groupID <= 0 {
		utils.HandleError(w, http.StatusBadRequest, "Invalid group id")
		return
	}

	if !services.IsGroupCreator(groupID, userID) {
		utils.HandleError(w, http.StatusForbidden, "Only the group creator can manage join requests")
		return
	}

	requests, err := repositories.ListRequestsForGroup(groupID)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to load requests")
		return
	}

	utils.RespondJSON(w, http.StatusOK, requests)
}

// RespondJoinRequestHandler lets the creator accept or refuse a join request.
func RespondJoinRequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID, err := auth.GetUserIDFromSession(r)
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, "You must be logged in")
		return
	}

	groupID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || groupID <= 0 {
		utils.HandleError(w, http.StatusBadRequest, "Invalid group id")
		return
	}

	requestID, err := strconv.Atoi(r.PathValue("requestID"))
	if err != nil || requestID <= 0 {
		utils.HandleError(w, http.StatusBadRequest, "Invalid request id")
		return
	}

	if !services.IsGroupCreator(groupID, userID) {
		utils.HandleError(w, http.StatusForbidden, "Only the group creator can manage join requests")
		return
	}

	action := r.PathValue("action")
	switch action {
	case "accept":
		err = repositories.AcceptRequest(requestID, userID)
	case "refuse":
		err = repositories.RefuseRequest(requestID)
	default:
		utils.HandleError(w, http.StatusBadRequest, "Invalid action, must be accept or refuse")
		return
	}

	if err != nil {
		utils.HandleError(w, http.StatusNotFound, "Request not found")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Request " + action + "d"})
}
