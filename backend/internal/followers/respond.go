package followers

import (
	"net/http"
	"strconv"

	"zone/internal/auth"
	"zone/internal/database"
	"zone/internal/utils"
)

func AcceptFollowAPI(w http.ResponseWriter, r *http.Request) {
	respondToRequest(w, r, "accepted")
}

func DeclineFollowAPI(w http.ResponseWriter, r *http.Request) {
	respondToRequest(w, r, "declined")
}

func respondToRequest(w http.ResponseWriter, r *http.Request, action string) {
	if r.Method != http.MethodPost {
		utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	followeeID, err := auth.GetUserID(r)
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, "You must be logged in")
		return
	}

	followerID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || followerID <= 0 {
		utils.HandleError(w, http.StatusBadRequest, "Invalid user id")
		return
	}

	if action == "accepted" {
		_, err = database.Database.Exec(
			"UPDATE followers SET status = 'accepted' WHERE follower_id = ? AND followee_id = ? AND status = 'pending'",
			followerID, followeeID,
		)
	} else {
		_, err = database.Database.Exec(
			"DELETE FROM followers WHERE follower_id = ? AND followee_id = ? AND status = 'pending'",
			followerID, followeeID,
		)
	}
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to update request")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Done"})
}