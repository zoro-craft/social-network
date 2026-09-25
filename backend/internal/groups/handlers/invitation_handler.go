package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"zone/internal/auth"
	"zone/internal/database"
	"zone/internal/groups/models"
	"zone/internal/groups/repositories"
	"zone/internal/groups/services"
	notificationservice "zone/internal/notifications/services"
	"zone/internal/utils"
)


// the handler is used to create a group invitation and send a notification to the invited user.
func InvitationsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.HandleError(w, 405, http.StatusText(405))
		return
	}

	userID, err := auth.GetUserID(r)
	if err != nil {
		utils.HandleError(w, 401, http.StatusText(401))
		return
	}

	var invitation models.GroupInvitation
	// Decode the request body into the invitation struct
	err = json.NewDecoder(r.Body).Decode(&invitation)
	if err != nil {
		utils.HandleError(w, 400, http.StatusText(400))
		return
	}

	// The current user is the inviter
	invitation.InviterID = userID

	// Validate the invitation creation
	if !services.ValidCreateInvitation(invitation, userID) {
		utils.HandleError(w, 400, "invalid invitation")
		return
	}

	// Use a transaction so the invitation and its notification
	// are created atomically (all succeed or none is committed).
	tx, err := database.Database.Begin()
	if err != nil {
		utils.HandleError(w, 500, "failed to create invitation")
		return
	}
	defer tx.Rollback()

	// Save invitation
	invitationID, err := repositories.SaveInvitation(tx, invitation)
	if err != nil {
		if err.Error() == "already invited" {
			utils.HandleError(w, http.StatusConflict, "already invited")
			return
		}
		utils.HandleError(w, 500, "failed to create invitation")
		return
	}

	groupName, err := repositories.GetGroupNameByID(invitation.GroupID)
	if err != nil {
		utils.HandleError(w, 500, "failed to get group")
		return
	}
	
	message := repositories.GetUserFullName(userID) + " invited you to join " + groupName

	if err := tx.Commit(); err != nil {
		utils.HandleError(w, 500, "failed to create invitation")
		return
	}

		// Create notification for the invited user
	_, err = notificationservice.Notifier.Notify(
		invitation.InvitedUserID,
		"group_invitation",
		&invitationID,
		&userID,
		message,
	)
	if err != nil {
		fmt.Println("CreateNotification error:", err)
		utils.HandleError(w, 500, "failed to create notification")
		return
	}


	utils.RespondJSON(
		w,
		http.StatusCreated,
		map[string]string{
			"message": "invitation created successfully!",
		},
	)
}
