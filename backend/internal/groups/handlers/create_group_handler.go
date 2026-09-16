package handlers

import (
	"encoding/json"
	"net/http"

	"zone/internal/auth"
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
	// read the request body and decode it into the Groupe struct
	err = json.NewDecoder(r.Body).Decode(&groupe)
	if err != nil {
		utils.HandleError(w, 400, http.StatusText(400))
		return
	}

	groupe.Admin = user_id
	// validate the group input 
	if !services.ValidateGroup(groupe) {
		utils.HandleError(w, 400, http.StatusText(400))
		return
	}

	groupe_Id , err := repositories.SaveGroup(groupe)
	if err != nil {
		utils.HandleError(w, 500, http.StatusText(500))
		return
	}
	
	err  = repositories.AddGroupMember(groupe_Id, user_id) 
	if err != nil {
		utils.HandleError(w, 500, http.StatusText(500))
		return
	}

	utils.RespondJSON(w, http.StatusCreated, map[string]string{"message": " Groupe created successfully!"})
}
