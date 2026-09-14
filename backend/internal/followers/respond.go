package followers

import (
	"net/http"
	"strconv"

	"zone/internal/auth"
	"zone/internal/database"
	"zone/internal/utils"
)

type FollowRequest struct {
	FollowerID int    `json:"follower_id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Avatar     string `json:"avatar"`
}

func GetRequestsAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID, err := auth.GetUserID(r)
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	rows, err := database.Database.Query(`
		SELECT f.follower_id, u.first_name, u.last_name, COALESCE(u.avatar, '')
		FROM followers f
		JOIN users u ON u.id = f.follower_id
		WHERE f.followee_id = ? AND f.status = 'pending'`, userID)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Database error")
		return
	}
	defer rows.Close()

	requests := make([]FollowRequest, 0)
	for rows.Next() {
		var req FollowRequest
		if err := rows.Scan(&req.FollowerID, &req.FirstName, &req.LastName, &req.Avatar); err != nil {
			utils.HandleError(w, http.StatusInternalServerError, "Database error")
			return
		}
		requests = append(requests, req)
	}
	if err := rows.Err(); err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Database error")
		return
	}

	utils.RespondJSON(w, http.StatusOK, requests)
}

func RespondRequestAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	followerID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || followerID <= 0 {
		utils.HandleError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	followeeID, err := auth.GetUserID(r)
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	switch action := r.PathValue("action"); action {
	case "accept":
		_, err = database.Database.Exec(
			"UPDATE followers SET status = 'accepted' WHERE follower_id = ? AND followee_id = ? AND status = 'pending'",
			followerID, followeeID,
		)
	case "decline":
		_, err = database.Database.Exec(
			"DELETE FROM followers WHERE follower_id = ? AND followee_id = ? AND status = 'pending'",
			followerID, followeeID,
		)
	default:
		utils.HandleError(w, http.StatusBadRequest, "Invalid action")
		return
	}
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to update request")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"status": "success"})
}
