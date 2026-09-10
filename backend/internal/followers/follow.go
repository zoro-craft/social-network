package followers

import (
	"net/http"
	"strconv"

	"zone/internal/auth"
	"zone/internal/database"
	"zone/internal/utils"
)

func UnfollowAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	followerID, err := auth.GetUserID(r)
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, "You must be logged in")
		return
	}

	followeeID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || followeeID <= 0 {
		utils.HandleError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	_, err = database.Database.Exec(
		"DELETE FROM followers WHERE follower_id = ? AND followee_id = ?",
		followerID, followeeID,
	)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to unfollow")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Unfollowed"})
}

func ToggleFollowAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	followeeID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || followeeID <= 0 {
		utils.HandleError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	followerID, err := auth.GetUserID(r)
	if err != nil || followerID == 0 {
		utils.HandleError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if followerID == followeeID {
		utils.HandleError(w, http.StatusBadRequest, "Cannot follow yourself")
		return
	}

	var isPrivate bool
	err = database.Database.QueryRow("SELECT is_private FROM users WHERE id = ?", followeeID).Scan(&isPrivate)
	if err != nil {
		utils.HandleError(w, http.StatusNotFound, "User not found")
		return
	}

	newStatus := "accepted"
	if isPrivate {
		newStatus = "pending"
	}

	// Atomic toggle: INSERT OR IGNORE inserts the relation when absent.
	// When it is already present the insert is ignored (RowsAffected == 0)
	// and the relation is removed instead. No read-then-write race, so a
	// concurrent double-click can never hit the unique constraint.
	res, err := database.Database.Exec(
		"INSERT OR IGNORE INTO followers (follower_id, followee_id, status) VALUES (?, ?, ?)",
		followerID, followeeID, newStatus,
	)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to toggle follow")
		return
	}

	affected, err := res.RowsAffected()
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to toggle follow")
		return
	}
	if affected == 0 {
		_, err = database.Database.Exec("DELETE FROM followers WHERE follower_id = ? AND followee_id = ?", followerID, followeeID)
		if err != nil {
			utils.HandleError(w, http.StatusInternalServerError, "Failed to unfollow")
			return
		}
		newStatus = "none"
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{
		"status": newStatus,
	})
}
