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

	// Use a transaction so adding the member and deleting the invitation are atomic.
	tx, err := database.Database.Begin()
	if err != nil {
		utils.HandleError(w, 500, http.StatusText(500))
		return
	}
	defer tx.Rollback()
	// check if the user is member 
	if !repositories.IsMember(user_id, invitation.GroupID) {
        err = repositories.AddGroupMember(tx, invitation.GroupID, user_id)
        if err != nil {
            utils.HandleError(w, 500, http.StatusText(500))
            return
        }
    }

	err = repositories.DeleteInvitation(tx, invitation.ID)
	if err != nil {
		utils.HandleError(w, 500, http.StatusText(500))
		return
	}

	joinRequestID, errJoinReq := repositories.GetJoinRequestIDByGroupAndUser(invitation.GroupID, user_id)
	if errJoinReq == nil {
		err = repositories.DeleteRequest(tx, joinRequestID)
		if err != nil {
			utils.HandleError(w, 500, http.StatusText(500))
			return
		}
	}

	if err := tx.Commit(); err != nil {
		utils.HandleError(w, 500, http.StatusText(500))
		return
	}

	_ = notificationservice.Notifier.MarkReadOnAction(user_id, "group_invitation", invitation_id)

	if errJoinReq == nil {
		creatorID, creatorErr := repositories.GetGroupCreator(invitation.GroupID)
		if creatorErr == nil {
			_ = notificationservice.Notifier.RemoveNotificationOnAction(creatorID, "join_request", joinRequestID)
		}
	}

	ws.DefaultHub.Subscribe(user_id, invitation.GroupID) // subscribe the new member to the group room

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

	// Use a transaction so deleting the invitation and marking it read are atomic.
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

	if err := tx.Commit(); err != nil {
		utils.HandleError(w, 500, http.StatusText(500))
		return
	}

	// Resolving the invitation marks its notification read.
	_ = notificationservice.Notifier.MarkReadOnAction(user_id, "group_invitation", invitation_id)

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": " invitation reject successfully!"})
}
