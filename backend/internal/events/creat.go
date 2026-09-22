package events

import (
	"encoding/json"
	"net/http"

	"zone/internal/auth"
	"zone/internal/database"
	"zone/internal/groups/repositories"
	notificationservice "zone/internal/notifications/services"
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
	if len(req.Title) >= 10 {
		utils.HandleError(w, http.StatusBadRequest, "title is to long")
		return
	}
	if req.Description == "" {
		utils.HandleError(w, http.StatusBadRequest, "Description is required")
		return
	}
	if len(req.Description) >= 30 {
		utils.HandleError(w, http.StatusBadRequest, "Description is to long")
		return
	}
	if req.EventDate == "" {
		utils.HandleError(w, http.StatusBadRequest, "event date is required")
		return
	}
	if req.EventTime == "" {
		utils.HandleError(w, http.StatusBadRequest, "event time is required")
		return
	}
	if req.FitsOptions == "" || req.SecondOptions == "" {
		utils.HandleError(w, http.StatusBadRequest, "both options are required")
		return
	}
	if len(req.FitsOptions) >= 10 || len(req.SecondOptions) >= 10 {
		utils.HandleError(w, http.StatusBadRequest, "both options is to long")
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
		"INSERT INTO events (group_id, title, description, event_date, event_time, created_by, options) VALUES (?, ?, ?, ?, ?, ?, ?)",
		req.GroupID, req.Title, req.Description, req.EventDate, req.EventTime, userID, string(optionsJSON),
	)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Could not create event")
		return
	}

	id, _ := res.LastInsertId()
	eventID := int(id)

	// Fan out an event_created notification to every group member except the creator.
	_ = notificationservice.Notifier.NotifyGroupEvent(req.GroupID, eventID, userID, req.Title, repositories.GetUserFullName(userID))

	utils.RespondJSON(w, http.StatusCreated, Event{
		ID:          eventID,
		GroupID:     req.GroupID,
		Title:       req.Title,
		Description: req.Description,
		EventDate:   req.EventDate,
		EventTime:   req.EventTime,
		CreatedBy:   userID,
	})
}
