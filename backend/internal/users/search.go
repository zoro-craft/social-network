package users

import (
	"net/http"
	"strings"

	"zone/internal/auth"
	"zone/internal/database"
	"zone/internal/utils"
)

func SearchUsersAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	viewerID, err := auth.GetUserID(r)
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, "You must be logged in to post")
		return
	}
	
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	if q == "" {
		utils.RespondJSON(w, http.StatusOK, []UserResult{})
		return
	}

	like := "%" + utils.EscapeLikeQuery(q) + "%"

	rows, err := database.Database.Query(`
		SELECT id, full_name, COALESCE(avatar, '')
		FROM users
		WHERE full_name LIKE ? AND id != ?
	`, like, viewerID)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Search failed")
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
    utils.HandleError(w, http.StatusInternalServerError, "Search failed")
    return
}
	utils.RespondJSON(w, http.StatusOK, results)
}