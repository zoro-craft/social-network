package followers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"

	"zone/internal/auth"
	"zone/internal/database"
	"zone/internal/groups/repositories"
	notificationservice "zone/internal/notifications/services"
	"zone/internal/utils"
)

var errUserNotFound = errors.New("follow target does not exist")

func ToggleFollowAPI(w http.ResponseWriter, r *http.Request) {
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

	if followerID == followeeID {
		utils.HandleError(w, http.StatusBadRequest, "Cannot follow yourself")
		return
	}

	newStatus, err := toggleFollow(followerID, followeeID)
	if err != nil {
		if errors.Is(err, errUserNotFound) {
			utils.HandleError(w, http.StatusNotFound, "User not found")
			return
		}
		log.Println("toggle follow:", err)
		utils.HandleError(w, http.StatusInternalServerError, "Failed to update follow status")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"status": newStatus})
}

func toggleFollow(followerID, followeeID int) (string, error) {
	tx, err := database.Database.Begin()
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	var isPrivate bool
	var status sql.NullString
	err = tx.QueryRow(`
        SELECT COALESCE(u.is_private, 0), f.status
        FROM users u
        LEFT JOIN followers f
            ON f.followee_id = u.id AND f.follower_id = ?
        WHERE u.id = ?`,
		followerID, followeeID,
	).Scan(&isPrivate, &status)
	if errors.Is(err, sql.ErrNoRows) {
		return "", errUserNotFound
	}
	if err != nil {
		return "", err
	}

	var newStatus string
	switch {
	case !status.Valid:
		newStatus = "accepted"
		if isPrivate {
			newStatus = "pending"
		}
		_, err = tx.Exec(
			"INSERT INTO followers (follower_id, followee_id, status) VALUES (?, ?, ?)",
			followerID, followeeID, newStatus,
		)

	case status.String == "accepted", isPrivate:
		newStatus = "none"
		_, err = tx.Exec(
			"DELETE FROM followers WHERE follower_id = ? AND followee_id = ?",
			followerID, followeeID,
		)

	default:
		newStatus = "accepted"
		_, err = tx.Exec(
			"UPDATE followers SET status = 'accepted' WHERE follower_id = ? AND followee_id = ? AND status = 'pending'",
			followerID, followeeID,
		)
	}
	if err != nil {
		return "", err
	}

if err := tx.Commit(); err != nil {
		return "", err
	}

	if newStatus == "pending" {
		message := "New follow request from " + repositories.GetUserFullName(followerID)
		_, _ = notificationservice.Notifier.NotifyFollowRequest(followeeID, followerID, message)
	}

	if newStatus == "none" && status.String == "pending" {
		_ = notificationservice.Notifier.RemoveFollowRequestNotification(followeeID, followerID)
	}

	return newStatus, nil
}
