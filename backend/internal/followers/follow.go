package followers

import (
	"net/http"
	"strconv"

	"zone/internal/auth"
	"zone/internal/database"
	"zone/internal/utils"
)

func FollowAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	followerID, err := auth.GetUserIDFromSession(r)
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, "You must be logged in")
		return
	}

	followeeID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || followeeID <= 0 {
		utils.HandleError(w, http.StatusBadRequest, "Invalid user id")
		return
	}

	if followeeID == followerID {
		utils.HandleError(w, http.StatusBadRequest, "You cannot follow yourself")
		return
	}

	var isPublic bool
	err = database.Database.QueryRow(
		"SELECT is_public FROM users WHERE id = ?", followeeID,
	).Scan(&isPublic)
	if err != nil {
		utils.HandleError(w, http.StatusNotFound, "User not found")
		return
	}

	status := "pending"
	if isPublic {
		status = "accepted"
	}

	_, err = database.Database.Exec(
		"INSERT OR IGNORE INTO followers (follower_id, followee_id, status) VALUES (?, ?, ?)",
		followerID, followeeID, status,
	)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to follow")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"status": status})
}

func UnfollowAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	followerID, err := auth.GetUserIDFromSession(r)
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, "You must be logged in")
		return
	}

	followeeID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || followeeID <= 0 {
		utils.HandleError(w, http.StatusBadRequest, "Invalid user id")
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