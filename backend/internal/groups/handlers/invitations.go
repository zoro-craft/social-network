package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"zone/internal/auth"
	"zone/internal/groups/repositories"
	"zone/internal/groups/services"
	"zone/internal/utils"
)

type inviteRequest struct {
	UserID int `json:"user_id"`
}

// InviteUserHandler lets a member invite a user to join the group.
func InviteUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	inviterID, err := auth.GetUserIDFromSession(r)
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

	if !services.IsGroupMember(groupID, inviterID) && !services.IsGroupCreator(groupID, inviterID) {
		utils.HandleError(w, http.StatusForbidden, "Only group members can invite users")
		return
	}

	var req inviteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.UserID <= 0 {
		utils.HandleError(w, http.StatusBadRequest, "Invalid user_id")
		return
	}

	if req.UserID == inviterID {
		utils.HandleError(w, http.StatusBadRequest, "You cannot invite yourself")
		return
	}

	if services.IsGroupMember(groupID, req.UserID) || services.IsGroupCreator(groupID, req.UserID) {
		utils.HandleError(w, http.StatusBadRequest, "User is already a member")
		return
	}

	already, err := repositories.HasInvitation(groupID, req.UserID)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to check invitation")
		return
	}
	if already {
		utils.HandleError(w, http.StatusBadRequest, "User already has a pending invitation")
		return
	}

	if err := repositories.CreateInvitation(groupID, inviterID, req.UserID); err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to invite user")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Invitation sent"})
}

// MyInvitationsHandler returns pending invitations for the logged-in user.
func MyInvitationsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID, err := auth.GetUserIDFromSession(r)
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, "You must be logged in")
		return
	}

	invitations, err := repositories.ListInvitationsForUser(userID)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to load invitations")
		return
	}

	utils.RespondJSON(w, http.StatusOK, invitations)
}

// RespondInvitationHandler lets the invited user accept or reject an invitation.
func RespondInvitationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID, err := auth.GetUserIDFromSession(r)
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, "You must be logged in")
		return
	}

	invitationID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || invitationID <= 0 {
		utils.HandleError(w, http.StatusBadRequest, "Invalid invitation id")
		return
	}

	action := r.PathValue("action")
	switch action {
	case "accept":
		err = repositories.AcceptInvitation(invitationID, userID)
	case "reject":
		err = repositories.RejectInvitation(invitationID, userID)
	default:
		utils.HandleError(w, http.StatusBadRequest, "Invalid action, must be accept or reject")
		return
	}

	if err != nil {
		utils.HandleError(w, http.StatusNotFound, "Invitation not found")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Invitation " + action + "ed"})
}
