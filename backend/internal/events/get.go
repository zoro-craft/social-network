package events

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"zone/internal/auth"
	"zone/internal/database"
	"zone/internal/utils"
)

func GetEvents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID, err := auth.GetUserID(r)
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	groupID, err := strconv.Atoi(r.URL.Query().Get("group_id"))
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Invalid group_id")
		return
	}

	rows, err := database.Database.Query(`
		SELECT id, group_id, title, created_by, created_at, options
		FROM events WHERE group_id = ? ORDER BY id DESC`,
		groupID,
	)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Could not load events")
		return
	}
	defer rows.Close()

	events := []Event{}
	for rows.Next() {
		var e Event
		var optionsJSON string
		if err := rows.Scan(&e.ID, &e.GroupID, &e.Title, &e.CreatedBy, &e.CreatedAt, &optionsJSON); err != nil {
			utils.HandleError(w, http.StatusInternalServerError, "Could not load events")
			return
		}

		options, err := parseOptions(optionsJSON)
		if err != nil {
			utils.HandleError(w, http.StatusInternalServerError, "Could not parse event options")
			return
		}
		e.Options = options

		voteCounts, userResponse, err := loadVoteSummary(e.ID, userID)
		if err != nil {
			utils.HandleError(w, http.StatusInternalServerError, "Could not load event votes")
			return
		}
		e.VoteCounts = voteCounts
		e.UserResponse = userResponse

		events = append(events, e)
	}

	utils.RespondJSON(w, http.StatusOK, events)
}

func parseOptions(optionsJSON string) ([]string, error) {
	var options []string
	err := json.Unmarshal([]byte(optionsJSON), &options)
	return options, err
}


func loadVoteSummary(eventID int, userID int) ([]VoteCount, string, error) {
	voteCounts := []VoteCount{}

	countRows, err := database.Database.Query(
		"SELECT response, COUNT(*) FROM event_responses WHERE event_id = ? GROUP BY response",
		eventID,
	)
	if err != nil {
		return nil, "", err
	}
	defer countRows.Close()

	for countRows.Next() {
		var vc VoteCount
		if err := countRows.Scan(&vc.Option, &vc.Count); err != nil {
			return nil, "", err
		}
		voteCounts = append(voteCounts, vc)
	}

	var userResponse string
	err = database.Database.QueryRow(
		"SELECT response FROM event_responses WHERE event_id = ? AND user_id = ?",
		eventID, userID,
	).Scan(&userResponse)
	if err != nil && err != sql.ErrNoRows {
		return nil, "", err
	}

	return voteCounts, userResponse, nil
}