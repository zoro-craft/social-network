package events

import (
	"encoding/json"
	"net/http"

	"zone/internal/auth"
	"zone/internal/database"
	"zone/internal/utils"
)

func CreateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID, err := auth.GetUserID(r)
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req createEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Title == "" {
		utils.HandleError(w, http.StatusBadRequest, "title is required")
		return
	}
	if len(req.Title)<=20 {
		utils.HandleError(w, http.StatusBadRequest, "title is to long")
		return
	}
	if req.FitsOptions == "" || req.SecondOptions == "" {
		utils.HandleError(w, http.StatusBadRequest, "both options are required")
		return
	}
	if len(req.FitsOptions) <=10|| len(req.SecondOptions) <=10 {
		utils.HandleError(w, http.StatusBadRequest, "both options are required")
		return
	}
	if !utils.IsGroupMember(req.GroupID, userID) {
		utils.HandleError(w, http.StatusForbidden, "You must be a group member to create an event")
		return
	}

	options := []string{req.FitsOptions, req.SecondOptions}
	optionsJSON, err := json.Marshal(options)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Invalid options")
		return
	}

	res, err := database.Database.Exec(
		"INSERT INTO events (group_id, title, created_by, options) VALUES (?, ?, ?, ?)",
		req.GroupID, req.Title, userID, string(optionsJSON),
	)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Could not create event")
		return
	}

	id, _ := res.LastInsertId()

	utils.RespondJSON(w, http.StatusCreated, Event{
		ID:        int(id),
		GroupID:   req.GroupID,
		Title:     req.Title,
		CreatedBy: userID,
	})
}
