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
	invitation_id, err := strconv.Atoi(parts[2])
	if err != nil {
		utils.HandleError(w, 400, http.StatusText(400))
		return
	}
	exists := repositories.InvitationExists(invitation_id , user_id) 
	if (!exists) {
		utils.HandleError(w, 400, http.StatusText(400))
		return
	}

}

func RejectInvitationResponse(w http.ResponseWriter, r *http.Request) {
}
