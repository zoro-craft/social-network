package handlers

import (
	"encoding/json"
	"net/http"

	"zone/internal/auth"
	"zone/internal/database"
	"zone/internal/groups/models"
	"zone/internal/groups/repositories"
	"zone/internal/groups/services"
	ws "zone/internal/websocket"

	"zone/internal/utils"
)


func CreatGroupeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.HandleError(w, 405, http.StatusText(405))
		return
	}
	user_id, err := auth.GetUserID(r)
	if err != nil {
		utils.HandleError(w, 401, http.StatusText(401))
		return
	}

	var groupe models.Groupe
	err = json.NewDecoder(r.Body).Decode(&groupe)
	if err != nil {
		utils.HandleError(w, 400, http.StatusText(400))
		return
	}
	groupe.Admin = user_id

	// Validate the group data
	if !services.ValidateGroup(groupe) {
		utils.HandleError(w, 400, http.StatusText(400))
		return
	}
	// start a transaction
	tx, err := database.Database.Begin()
	if err != nil {
		utils.HandleError(w, 500, http.StatusText(500))
		return
	}
	defer tx.Rollback()

	groupe_Id, err := repositories.SaveGroup(tx, groupe)
	if err != nil {
		utils.HandleError(w, 500, http.StatusText(500))
		return
	}
	ws.DefaultHub.Subscribe(user_id, groupe_Id) // subscribe creator to the group room (FR-007)

	err = repositories.AddGroupMember(tx, groupe_Id, user_id)
	if err != nil {
		utils.HandleError(w, 500, http.StatusText(500))
		return
	}
	// Commit to permanently save the changes to the database.
	if err := tx.Commit(); err != nil {
		utils.HandleError(w, 500, http.StatusText(500))
		return
	}

	utils.RespondJSON(w, http.StatusCreated, map[string]string{"message": " Groupe created successfully!"})
}
