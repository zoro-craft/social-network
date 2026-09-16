package handlers

import (
	"encoding/json"
	"net/http"

	"zone/internal/auth"
	"zone/internal/database"
	"zone/internal/groups/models"
	"zone/internal/groups/repositories"
	"zone/internal/groups/services"

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

	if !services.ValidateGroup(groupe) {
		utils.HandleError(w, 400, http.StatusText(400))
		return
	}

	// Use a transaction so the group and its creator member
	// are inserted atomically (all succeed or none is committed).
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

	err = repositories.AddGroupMember(tx, groupe_Id, user_id)
	if err != nil {
		utils.HandleError(w, 500, http.StatusText(500))
		return
	}

	if err := tx.Commit(); err != nil {
		utils.HandleError(w, 500, http.StatusText(500))
		return
	}

	utils.RespondJSON(w, http.StatusCreated, map[string]string{"message": " Groupe created successfully!"})
}
