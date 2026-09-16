package followers

import (
	"net/http"
	"strconv"

	"zone/internal/auth"
	"zone/internal/database"
	"zone/internal/utils"
)

type UserItem struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
}

func canViewList(viewerID, profileID int) bool {
	if viewerID == profileID {
		return true
	}

	var isPrivate bool
	if err := database.Database.QueryRow(
		"SELECT is_private FROM users WHERE id = ?", profileID,
	).Scan(&isPrivate); err != nil {
		return false
	}
	if !isPrivate {
		return true
	}

	var isFollower bool
	err := database.Database.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM followers WHERE follower_id = ? AND followee_id = ? AND status = 'accepted')",
		viewerID, profileID,
	).Scan(&isFollower)
	return err == nil && isFollower
}

func GetFollowersAPI(w http.ResponseWriter, r *http.Request) {
	listHandler(w, r, true)
}

func GetFollowingAPI(w http.ResponseWriter, r *http.Request) {
	listHandler(w, r, false)
}

func listHandler(w http.ResponseWriter, r *http.Request, listFollowers bool) {
	profileID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || profileID <= 0 {
		utils.HandleError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	viewerID, _ := auth.GetUserID(r)
	if !canViewList(viewerID, profileID) {
		utils.RespondJSON(w, http.StatusOK, []UserItem{})
		return
	}

	listUsers(w, profileID, listFollowers)
}

func listUsers(w http.ResponseWriter, profileID int, listFollowers bool) {
	query := `
		SELECT u.id, u.first_name, u.last_name, COALESCE(u.nickname, ''), COALESCE(u.avatar, '')
		FROM followers f
		JOIN users u ON u.id = f.follower_id
		WHERE f.followee_id = ? AND f.status = 'accepted'`
	if !listFollowers {
		query = `
			SELECT u.id, u.first_name, u.last_name, COALESCE(u.nickname, ''), COALESCE(u.avatar, '')
			FROM followers f
			JOIN users u ON u.id = f.followee_id
			WHERE f.follower_id = ? AND f.status = 'accepted'`
	}

	rows, err := database.Database.Query(query, profileID)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Database error")
		return
	}
	defer rows.Close()

	users := make([]UserItem, 0)
	for rows.Next() {
		var u UserItem
		if err := rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Nickname, &u.Avatar); err != nil {
			utils.HandleError(w, http.StatusInternalServerError, "Database error")
			return
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Database error")
		return
	}

	utils.RespondJSON(w, http.StatusOK, users)
}
