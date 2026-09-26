package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"zone/internal/auth"
	"zone/internal/database"
	"zone/internal/events/models"
	"zone/internal/utils"
)

func VoteEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID, err := auth.GetUserID(r)
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req models.VoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	var optionsJSON string
	var groupID int
	err = database.Database.QueryRow(
		"SELECT group_id, options FROM events WHERE id = ?", req.EventID,
	).Scan(&groupID, &optionsJSON)
	if err == sql.ErrNoRows {
		utils.HandleError(w, http.StatusNotFound, "Event not found")
		return
	}
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Could not load event")
		return
	}
	if !utils.IsGroupMember(groupID, userID) {
		utils.HandleError(w, http.StatusForbidden, "You must be a group member to vote")
		return
	}

	options, err := parseOptions(optionsJSON)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Could not parse event options")
		return
	}

	valid := false
	for _, opt := range options {
		if opt == req.Vote {
			valid = true
			break
		}
	}
	if !valid {
		utils.HandleError(w, http.StatusBadRequest, "vote must match one of the event options")
		return
	}

	var existingResponse string
	err = database.Database.QueryRow(
		"SELECT response FROM event_responses WHERE user_id = ? AND event_id = ?",
		userID, req.EventID,
	).Scan(&existingResponse)

	switch {
	case err == sql.ErrNoRows:
		_, err = database.Database.Exec(
			"INSERT INTO event_responses (user_id, event_id, response) VALUES (?, ?, ?)",
			userID, req.EventID, req.Vote,
		)
	case err == nil && existingResponse == req.Vote:
		_, err = database.Database.Exec(
			"DELETE FROM event_responses WHERE user_id = ? AND event_id = ?",
			userID, req.EventID,
		)
	case err == nil:
		_, err = database.Database.Exec(
			"UPDATE event_responses SET response = ? WHERE user_id = ? AND event_id = ?",
			req.Vote, userID, req.EventID,
		)
	}
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Could not save vote")
		return
	}

	voteCounts, userResponse, err := loadVoteSummary(req.EventID, userID)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Could not load vote counts")
		return
	}

	e := models.Event{
		ID:           req.EventID,
		VoteCounts:   voteCounts,
		UserResponse: userResponse,
	}

	utils.RespondJSON(w, http.StatusOK, e)
}
