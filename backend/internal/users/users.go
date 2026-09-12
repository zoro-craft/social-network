package users

import (
	"net/http"
	"strconv"
	"strings"

	"zone/internal/auth"
	"zone/internal/database"
	"zone/internal/utils"
)

type User struct {
	ID        int    `json:"id"`
	Full_name string `json:"full_name"`
	Avatar    string `json:"avatar,omitempty"`
}

func GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	user_id, err := auth.GetUserID(r)
	if err != nil {
		utils.HandleError(w, 401, http.StatusText(401))
		return

	}
	users, err := GetAllUsers(user_id)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Could not load users")
		return
	}

	utils.RespondJSON(w, http.StatusOK, users)
}

func GetAllUsersNotmemberHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	user_id, err := auth.GetUserID(r)
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, http.StatusText(401))
		return
	}

	// Example: /groups/5/users
	parts := strings.Split(r.URL.Path, "/")

	if len(parts) < 4 {
		utils.HandleError(w, http.StatusBadRequest, "Invalid group ID")
		return
	}

	group_id, err := strconv.Atoi(parts[2])
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Invalid group ID")
		return
	}

	users, err := GetAllUsersNotmember(user_id, group_id)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Could not load users")
		return
	}

	utils.RespondJSON(w, http.StatusOK, users)
}

func GetAllUsers(userid int) ([]User, error) {
	rows, err := database.Database.Query(`
        SELECT id, full_name, COALESCE(avatar, '')
        FROM users
        WHERE id != ?
        ORDER BY id
    `, userid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var u User

		if err := rows.Scan(
			&u.ID,
			&u.Full_name,
			&u.Avatar,
		); err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	return users, rows.Err()
}

func GetAllUsersNotmember(userid int, groupid int) ([]User, error) {
	rows, err := database.Database.Query(`
		SELECT id, full_name, COALESCE(avatar, '')
		FROM users
		WHERE id != ?
		AND id NOT IN (
			SELECT user_id
			FROM group_members
			WHERE group_id = ?
		)
		ORDER BY id
	`, userid, groupid)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []User

	for rows.Next() {
		var u User

		if err := rows.Scan(
			&u.ID,
			&u.Full_name,
			&u.Avatar,
		); err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	return users, rows.Err()
}
