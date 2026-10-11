package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

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

	var group models.Groupe
	isMultipart := strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data")

	if isMultipart {
		r.Body = http.MaxBytesReader(w, r.Body, utils.MaxImageFormSize)
		if err := r.ParseMultipartForm(utils.MaxUploadSize); err != nil {
			utils.HandleError(w, 400, "Invalid form data or file too large")
			return
		}
		group.Title = strings.TrimSpace(r.FormValue("titel"))
		group.Description = strings.TrimSpace(r.FormValue("description"))
	} else {
		if err := json.NewDecoder(r.Body).Decode(&group); err != nil {
			utils.HandleError(w, 400, http.StatusText(400))
			return
		}
		group.Title = strings.TrimSpace(group.Title)
		group.Description = strings.TrimSpace(group.Description)
		group.Image = ""
	}
	group.Admin = user_id

	if !services.ValidateGroup(group) {
		utils.HandleError(w, 400, http.StatusText(400))
		return
	}

	if isMultipart {
		group.Image, err = utils.SaveImageField(r, "image", "groups", "/uploads/groups/")
		if err != nil {
			utils.HandleError(w, 400, err.Error())
			return
		}
	}

	// start a transaction
	tx, err := database.Database.Begin()
	if err != nil {
		utils.HandleError(w, 500, http.StatusText(500))
		return
	}
	defer tx.Rollback()

	groupID, err := repositories.SaveGroup(tx, group)
	if err != nil {
		utils.HandleError(w, 500, http.StatusText(500))
		return
	}
	ws.DefaultHub.Subscribe(user_id, groupID) // subscribe the creator to the group room

	err = repositories.AddGroupMember(tx, groupID, user_id)
	if err != nil {
		utils.HandleError(w, 500, http.StatusText(500))
		return
	}
	// Commit to permanently save the changes to the database.
	if err := tx.Commit(); err != nil {
		utils.HandleError(w, 500, http.StatusText(500))
		return
	}

	utils.RespondJSON(w, http.StatusCreated, map[string]any{"message": "Group created successfully!", "group_id": groupID})
}
