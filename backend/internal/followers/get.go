package followers

import (
	"net/http"
	"strconv"

	"zone/internal/auth"
	"zone/internal/database"
	"zone/internal/utils"
)

func GetFollowersAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || userID <= 0 {
		utils.HandleError(w, http.StatusBadRequest, "Invalid user id")
		return
	}

	rows, err := database.Database.Query(`
		SELECT u.id, COALESCE(u.first_name, ''), COALESCE(u.avatar, '')
		FROM followers f
		JOIN users u ON u.id = f.follower_id
		WHERE f.followee_id = ? AND f.status = 'accepted'
	`, userID)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Database error")
		return
	}
	defer rows.Close()

	result := []Follow{}
	for rows.Next() {
		var f Follow
		if err := rows.Scan(&f.FollowerID, &f.FirstName, &f.Avatar); err != nil {
			continue
		}
		result = append(result, f)
	}

	utils.RespondJSON(w, http.StatusOK, result)
}

func GetFollowingAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || userID <= 0 {
		utils.HandleError(w, http.StatusBadRequest, "Invalid user id")
		return
	}

	rows, err := database.Database.Query(`
		SELECT u.id, COALESCE(u.first_name, ''), COALESCE(u.avatar, '')
		FROM followers f
		JOIN users u ON u.id = f.followee_id
		WHERE f.follower_id = ? AND f.status = 'accepted'
	`, userID)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Database error")
		return
	}
	defer rows.Close()

	result := []Follow{}
	for rows.Next() {
		var f Follow
		if err := rows.Scan(&f.FolloweeID, &f.FirstName, &f.Avatar); err != nil {
			continue
		}
		result = append(result, f)
	}

	utils.RespondJSON(w, http.StatusOK, result)
}

func GetFollowStatusAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	viewerID, err := auth.GetUserIDFromSession(r)
	if err != nil {
		utils.RespondJSON(w, http.StatusOK, map[string]string{"status": "none"})
		return
	}

	targetID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || targetID <= 0 {
		utils.HandleError(w, http.StatusBadRequest, "Invalid user id")
		return
	}

	var status string
	err = database.Database.QueryRow(
		"SELECT status FROM followers WHERE follower_id = ? AND followee_id = ?",
		viewerID, targetID,
	).Scan(&status)
	if err != nil {
		status = "none"
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"status": status})
}