package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"zone/internal/auth"
	"zone/internal/groups/models"
	"zone/internal/groups/repositories"
	"zone/internal/groups/services"
	"zone/internal/utils"
)

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

	// Save invitation
	invitationID, err := repositories.SaveInvitation(invitation)
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

	message := "You are invited to join " + groupName

	// Create notification for the invited user
	err = repositories.CreateNotification(
		invitation.InvitedUserID,
		"group_invtation",
		&invitationID,
		nil,
		message,
	)
	if err != nil {
		fmt.Println("CreateNotification error:", err)
		utils.HandleError(w, 500, "failed to create notification")
		return
	}
	fmt.Println("2")

	utils.RespondJSON(
		w,
		http.StatusCreated,
		map[string]string{
			"message": "invitation created successfully!",
		},
	)
}
