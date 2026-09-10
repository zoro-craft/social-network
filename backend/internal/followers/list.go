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
	err := database.Database.QueryRow("SELECT is_private FROM users WHERE id = ?", profileID).Scan(&isPrivate)
	if err != nil {
		return false
	}
	if !isPrivate {
		return true
	}
	var status string
	err = database.Database.QueryRow("SELECT status FROM followers WHERE follower_id = ? AND followee_id = ?", viewerID, profileID).Scan(&status)
	return err == nil && status == "accepted"
}

func GetFollowersAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	profileID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Invalid ID")
		return
	}
	viewerID, _ := auth.GetUserID(r)

	if !canViewList(viewerID, profileID) {
		utils.RespondJSON(w, http.StatusOK, []UserItem{})
		return
	}

	rows, err := database.Database.Query(`
		SELECT u.id, u.first_name, u.last_name, COALESCE(u.nickname, ''), COALESCE(u.avatar, '') 
		FROM users u 
		JOIN followers f ON u.id = f.follower_id 
		WHERE f.followee_id = ? AND f.status = 'accepted'`, profileID)

	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Database error")
		return
	}
	defer rows.Close()

	var users = make([]UserItem, 0)
	for rows.Next() {
		var u UserItem
		rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Nickname, &u.Avatar)
		users = append(users, u)
	}
	utils.RespondJSON(w, http.StatusOK, users)
}

func GetFollowingAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	profileID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Invalid ID")
		return
	}
	viewerID, _ := auth.GetUserID(r)

	if !canViewList(viewerID, profileID) {
		utils.RespondJSON(w, http.StatusOK, []UserItem{})
		return
	}

	rows, err := database.Database.Query(`
		SELECT u.id, u.first_name, u.last_name, COALESCE(u.nickname, ''), COALESCE(u.avatar, '') 
		FROM users u 
		JOIN followers f ON u.id = f.followee_id 
		WHERE f.follower_id = ? AND f.status = 'accepted'`, profileID)

	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Database error")
		return
	}
	defer rows.Close()

	var users = make([]UserItem, 0)
	for rows.Next() {
		var u UserItem
		rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Nickname, &u.Avatar)
		users = append(users, u)
	}
	utils.RespondJSON(w, http.StatusOK, users)
}