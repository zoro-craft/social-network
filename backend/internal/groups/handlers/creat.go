package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"zone/internal/auth"
	"zone/internal/groups/models"
	"zone/internal/groups/repositories"
	"zone/internal/groups/services"
	"zone/internal/utils"
)

func CreateGroupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	userID, err := auth.GetUserIDFromSession(r)
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, "You must be logged in")
		return
	}

	var group models.Group
	err = json.NewDecoder(r.Body).Decode(&group)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}
	group.Creator = userID
	group.Title = strings.TrimSpace(group.Title)
	group.Description = strings.TrimSpace(group.Description)

	if !services.ValidGroup(group) {
		utils.HandleError(w, http.StatusBadRequest, "Invalid group: title and description are required (title < 50, description < 200)")
		return
	}

	id, err := repositories.CreateGroup(group.Title, group.Description, userID)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to create group")
		return
	}

	// The creator automatically becomes a member.
	if err := repositories.AddMember(int(id), userID); err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to add creator as member")
		return
	}

	group.ID = int(id)
	utils.RespondJSON(w, http.StatusCreated, group)
}

