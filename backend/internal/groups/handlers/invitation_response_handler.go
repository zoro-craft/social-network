package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"zone/internal/auth"
	"zone/internal/database"
	"zone/internal/groups/repositories"
	"zone/internal/utils"
	ws "zone/internal/websocket"
)

// POST /invitations/1/accept
// POST /invitations/1/reject

func AcceptInvitationResponse(w http.ResponseWriter, r *http.Request) {
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

	// Use a transaction so adding the member, updating the invitation
	// status and removing the notification are applied atomically.
	tx, err := database.Database.Begin()
	if err != nil {
		utils.HandleError(w, 500, http.StatusText(500))
		return
	}
	defer tx.Rollback()

	if (repositories.IsMember(user_id, invitation.GroupID)) {
		repositories.DeleteInvitation(tx , invitation.ID)
		repositories.DeleteInvitationNotification(tx, user_id, invitation_id)
		tx.Commit()
		return
	}

	err = repositories.AddGroupMember(tx, invitation.GroupID, user_id)
	if err != nil {
		utils.HandleError(w, 500, http.StatusText(500))
		return
	}

	// Delete the invitation when the user accept
	err = repositories.DeleteInvitation(tx, invitation.ID)
	if err != nil {
		utils.HandleError(w, 500, http.StatusText(500))
		return
	}

	// Remove the invitation notification
	err = repositories.DeleteInvitationNotification(tx, user_id, invitation_id)
	if err != nil {
		utils.HandleError(w, 500, http.StatusText(500))
		return
	}

	if err := tx.Commit(); err != nil {
		utils.HandleError(w, 500, http.StatusText(500))
		return
	}

	ws.DefaultHub.Subscribe(user_id, invitation.GroupID) // subscribe new member to the group room (FR-007)

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": " invitation accept successfully!"})
}

func RejectInvitationResponse(w http.ResponseWriter, r *http.Request) {
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

	// Use a transaction so updating the invitation status and
	// removing the notification are applied atomically.
	tx, err := database.Database.Begin()
	if err != nil {
		utils.HandleError(w, 500, http.StatusText(500))
		return
	}
	defer tx.Rollback()

	err = repositories.DeleteInvitation(tx, invitation.ID)
	if err != nil {
		utils.HandleError(w, 500, http.StatusText(500))
		return
	}

	err = repositories.DeleteInvitationNotification(tx, user_id, invitation_id)
	if err != nil {
		utils.HandleError(w, 500, http.StatusText(500))
		return
	}

	if err := tx.Commit(); err != nil {
		utils.HandleError(w, 500, http.StatusText(500))
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": " invitation reject successfully!"})
}
