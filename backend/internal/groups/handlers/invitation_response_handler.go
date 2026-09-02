package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"zone/internal/auth"
	"zone/internal/groups/repositories"
	"zone/internal/utils"
)

// POST /invitations/1/accept
// POST /invitations/1/reject

// accept ->  check_validation_of_invitation -> change stat -> instrt the user to group_memmber  -> remove the invitation from database
// reject -> check_validation_of_invitation -> change stat -> remove the invitation from datatbase

func AcceptInvitationResponse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.HandleError(w, 405, http.StatusText(405))
		return
	}
	user_id, err := auth.GetUserIDFromSession(r)
	if err != nil {
		utils.HandleError(w, 401, http.StatusText(401))
		return
	}
	parts := strings.Split((r.URL.Path), "/")

	if len(parts) < 4 {
		utils.HandleError(w, 400, http.StatusText(400))
		return
	}
	invitation_id, err := strconv.Atoi(parts[2])
	if err != nil {
		utils.HandleError(w, 400, http.StatusText(400))
		return
	}
	invitation, err := repositories.GetInvitation(invitation_id, user_id)
	if err != nil {
		utils.HandleError(w, 400, http.StatusText(400))
		return
	}
	if invitation.Status != "pending" {
		utils.HandleError(w, 400, http.StatusText(400))
		return
	}

	err = repositories.AddGroupMember(
		invitation.GroupID,
		user_id,
	)
	if err != nil {
		utils.HandleError(w, 500, http.StatusText(500))
		return
	}
	err = repositories.UpdateInvitationStatus(
		invitation_id,
		"accepted",
	)
	if err != nil {
		utils.HandleError(w, 500, http.StatusText(500))
		return
	}
	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": " invitation accept successfully!"})
}

func RejectInvitationResponse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.HandleError(w, 405, http.StatusText(405))
		return
	}
	user_id, err := auth.GetUserIDFromSession(r)
	if err != nil {
		utils.HandleError(w, 401, http.StatusText(401))
		return
	}
	parts := strings.Split((r.URL.Path), "/")

	if len(parts) < 4 {
		utils.HandleError(w, 400, http.StatusText(400))
		return
	}
	invitation_id, err := strconv.Atoi(parts[2])
	if err != nil {
		utils.HandleError(w, 400, http.StatusText(400))
		return
	}
	invitation, err := repositories.GetInvitation(invitation_id, user_id)
	if err != nil {
		utils.HandleError(w, 400, http.StatusText(400))
		return
	}
	if invitation.Status != "pending" {
		utils.HandleError(w, 400, http.StatusText(400))
		return
	}
	err = repositories.UpdateInvitationStatus(
		invitation_id,
		"rejected",
	)
	if err != nil {
		utils.HandleError(w, 500, http.StatusText(500))
		return
	}
	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": " invitation reject successfully!"})
}
