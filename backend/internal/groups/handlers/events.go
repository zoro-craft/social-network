package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"zone/internal/auth"
	"zone/internal/groups/models"
	"zone/internal/groups/repositories"
	"zone/internal/groups/services"
	"zone/internal/utils"
)

type createEventRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	EventTime   string `json:"event_time"`
}

// CreateEventHandler lets a member create an event for the group.
func CreateEventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	groupID, userID, ok := requireGroupMember(w, r)
	if !ok {
		return
	}

	var req createEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	req.Title = strings.TrimSpace(req.Title)
	req.Description = strings.TrimSpace(req.Description)
	req.EventTime = strings.TrimSpace(req.EventTime)

	if req.Title == "" {
		utils.HandleError(w, http.StatusBadRequest, "Title is required")
		return
	}
	if req.EventTime == "" {
		utils.HandleError(w, http.StatusBadRequest, "event_time is required")
		return
	}
	if _, err := time.Parse(time.RFC3339, req.EventTime); err != nil {
		utils.HandleError(w, http.StatusBadRequest, "event_time must be in RFC3339 format (e.g. 2026-01-01T18:00:00Z)")
		return
	}
	if len(req.Title) > 100 {
		utils.HandleError(w, http.StatusBadRequest, "Title too long")
		return
	}
	if len(req.Description) > 2000 {
		utils.HandleError(w, http.StatusBadRequest, "Description too long")
		return
	}

	event := models.Event{
		GroupID:     groupID,
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		EventTime:   req.EventTime,
	}

	id, err := repositories.CreateEvent(event)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to create event")
		return
	}

	event.ID = int(id)
	utils.RespondJSON(w, http.StatusCreated, event)
}

// EventsListHandler returns all events for a group (members only).
func EventsListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	groupID, userID, ok := requireGroupMember(w, r)
	if !ok {
		return
	}

	events, err := repositories.ListEvents(groupID, userID)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to load events")
		return
	}

	utils.RespondJSON(w, http.StatusOK, events)
}

type respondEventRequest struct {
	Response string `json:"response"`
}

// RespondEventHandler lets a member set their participation for an event.
func RespondEventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID, err := auth.GetUserIDFromSession(r)
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, "You must be logged in")
		return
	}

	eventID, err := strconv.Atoi(r.PathValue("eventID"))
	if err != nil || eventID <= 0 {
		utils.HandleError(w, http.StatusBadRequest, "Invalid event id")
		return
	}

	event, err := repositories.GetEvent(eventID)
	if err != nil {
		utils.HandleError(w, http.StatusNotFound, "Event not found")
		return
	}

	if !services.IsGroupMember(event.GroupID, userID) && !services.IsGroupCreator(event.GroupID, userID) {
		utils.HandleError(w, http.StatusForbidden, "You must be a group member")
		return
	}

	var req respondEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	switch req.Response {
	case "going", "not_going", "maybe":
	default:
		utils.HandleError(w, http.StatusBadRequest, "response must be going, not_going, or maybe")
		return
	}

	if err := repositories.SetEventResponse(eventID, userID, req.Response); err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to save response")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"response": req.Response})
}
