package handlers

import (
	"net/http"

	"zone/internal/auth"
	"zone/internal/groups/repositories"
	"zone/internal/utils"
)

func GetInvitationsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.HandleError(w, 405, http.StatusText(405))
		return
	}

	userID, err := auth.GetUserIDFromSession(r)
	if err != nil {
		utils.HandleError(w, 401, http.StatusText(401))
		return
	}

	invitations, err := repositories.GetInvitationsByUserID(userID)
	if err != nil {
		utils.HandleError(w, 500, "failed to get invitations")
		return
	}

	utils.RespondJSON(w, http.StatusOK, invitations)
}
