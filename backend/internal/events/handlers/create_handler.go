package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"zone/internal/auth"
	"zone/internal/database"
	"zone/internal/events/models"
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

	var req models.CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Title == "" {
		utils.HandleError(w, http.StatusBadRequest, "title is required")
		return
	}
	if len(req.Title) > 50 {
		utils.HandleError(w, http.StatusBadRequest, "title is too long")
		return
	}
	if req.Description == "" {
		utils.HandleError(w, http.StatusBadRequest, "Description is required")
		return
	}
	if len(req.Description) > 100 {
		utils.HandleError(w, http.StatusBadRequest, "Description is too long")
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

	if _, err := time.Parse("15:04", req.EventTime); err != nil {
		utils.HandleError(w, http.StatusBadRequest, "invalid event time format")
		return
	}
	if _, err := time.Parse("2006-01-02", req.EventDate); err != nil {
		utils.HandleError(w, http.StatusBadRequest, "invalid event date format")
		return
	}

	maxDate := time.Now().AddDate(10, 0, 0).Format("2006-01-02")
	nowTime := time.Now().Format("15:04")
	today := time.Now().Format("2006-01-02")

	if req.EventDate == today && req.EventTime < nowTime {
		utils.HandleError(w, http.StatusBadRequest, "event time cannot be in the past")
		return
	}
	if req.EventDate > maxDate {
		utils.HandleError(w, http.StatusBadRequest, "event date cannot be more than 10 years in the future")
		return
	}

	if req.EventDate < today {
		utils.HandleError(w, http.StatusBadRequest, "event date cannot be in the past")
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
	if req.FitsOptions == req.SecondOptions {
		utils.HandleError(w, http.StatusBadRequest, "Both participation options are required")
		return
	}
	if len(req.FitsOptions) > 20 || len(req.SecondOptions) > 20 {
		utils.HandleError(w, http.StatusBadRequest, "both options are too long")
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

	utils.RespondJSON(w, http.StatusCreated, models.Event{
		ID:          eventID,
		GroupID:     req.GroupID,
		Title:       req.Title,
		Description: req.Description,
		EventDate:   req.EventDate,
		EventTime:   req.EventTime,
		CreatedBy:   userID,
	})
}
