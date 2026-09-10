package followers

import (
	"net/http"
	"strconv"
	"zone/internal/auth"
	"zone/internal/database"
	"zone/internal/utils"
)

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
		SELECT u.id, u.first_name, u.last_name, COALESCE(u.avatar, '')
		FROM users u
		JOIN followers f ON u.id = f.follower_id
		WHERE f.followee_id = ? AND f.status = 'pending'`, userID)

	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Database error")
		return
	}
	defer rows.Close()

	var requests = make([]map[string]interface{}, 0)
	for rows.Next() {
		var id int
		var first, last, avatar string
		rows.Scan(&id, &first, &last, &avatar)
		requests = append(requests, map[string]interface{}{
			"follower_id": id,
			"first_name":  first,
			"last_name":   last,
			"avatar":      avatar,
		})
	}
	utils.RespondJSON(w, http.StatusOK, requests)
}

func RespondRequestAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	followerID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Invalid ID")
		return
	}
	action := r.PathValue("action")
	userID, _ := auth.GetUserID(r)

	if action == "accept" {
		_, err = database.Database.Exec("UPDATE followers SET status = 'accepted' WHERE follower_id = ? AND followee_id = ?", followerID, userID)
	} else if action == "decline" {
		_, err = database.Database.Exec("DELETE FROM followers WHERE follower_id = ? AND followee_id = ?", followerID, userID)
	} else {
		utils.HandleError(w, http.StatusBadRequest, "Invalid action")
		return
	}

	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to update request")
		return
	}
	utils.RespondJSON(w, http.StatusOK, map[string]string{"status": "success"})
}

