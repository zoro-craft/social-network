package users

import (
	"net/http"
	"strings"

	"zone/internal/auth"
	"zone/internal/database"
	"zone/internal/utils"
	"zone/internal/websocket"
)

// OnlineUsersAPI serves GET /users/online: every connected user except the
// viewer, as id/full_name/avatar rows for the sidebar.
func OnlineUsersAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	viewerID, err := auth.GetUserID(r)
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, "You must be logged in")
		return
	}

	ids := websocket.DefaultHub.OnlineUserIDs()
	filtered := ids[:0]
	for _, id := range ids {
		if id != viewerID {
			filtered = append(filtered, id)
		}
	}
	if len(filtered) == 0 {
		utils.RespondJSON(w, http.StatusOK, []UserResult{})
		return
	}

	// The ids come from the hub, not the client, so they are safe. The IN list
	// still uses placeholders rather than string-built values.
	placeholders := make([]string, len(filtered))
	args := make([]any, len(filtered))
	for i, id := range filtered {
		placeholders[i] = "?"
		args[i] = id
	}

	rows, err := database.Database.Query(
		"SELECT id, full_name, COALESCE(avatar, '') FROM users WHERE id IN ("+
			strings.Join(placeholders, ",")+") ORDER BY full_name",
		args...,
	)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to load online users")
		return
	}
	defer rows.Close()

	results := []UserResult{}
	for rows.Next() {
		var u UserResult
		if err := rows.Scan(&u.ID, &u.FullName, &u.Avatar); err != nil {
			continue
		}
		results = append(results, u)
	}
	if err := rows.Err(); err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to load online users")
		return
	}

	utils.RespondJSON(w, http.StatusOK, results)
}
