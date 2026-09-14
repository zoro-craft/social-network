package comments

import (
	"net/http"
	"strconv"
	"zone/internal/auth"
	"zone/internal/database"
	"zone/internal/utils"
)

func GetCommentsByGroupAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID, err := auth.GetUserID(r)
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, "You must be logged in")
		return
	}

	groupID, err := strconv.Atoi(r.PathValue("group_id"))
	if err != nil || groupID <= 0 {
		utils.HandleError(w, http.StatusBadRequest, "Invalid group id")
		return
	}

	if !utils.IsGroupMember(groupID, userID) {
		utils.HandleError(w, http.StatusForbidden, "You must be a group member")
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 50 {
		limit = 10
	}
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if offset < 0 {
		offset = 0
	}

	beforeStr := r.URL.Query().Get("before")
	var beforeID int
	if beforeStr != "" {
		beforeID, err = strconv.Atoi(beforeStr)
		if err != nil || beforeID <= 0 {
			utils.HandleError(w, http.StatusBadRequest, "Invalid before id")
			return
		}
	}

	query := commentSelect + `
		WHERE c.group_id = ?
	`
	args := []interface{}{groupID}

	if beforeStr != "" {
		query += " AND c.id <= ?"
		args = append(args, beforeID)
	}

	query += " ORDER BY c.id DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := database.Database.Query(query, args...)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Database error loading comments")
		return
	}
	defer rows.Close()

	result := []Comment{}
	for rows.Next() {
		var c Comment
		if err := rows.Scan(&c.ID, &c.PostID, &c.GroupID, &c.Content, &c.Image, &c.UserID, &c.FirstName, &c.Avatar, &c.CreatedAt); err == nil {
			result = append(result, c)
		}
	}

	if beforeStr == "" {
		var maxID int
		database.Database.QueryRow(
			"SELECT COALESCE(MAX(id), 0) FROM comments WHERE group_id = ?", groupID,
		).Scan(&maxID)
		w.Header().Set("X-Max-Comment-Id", strconv.Itoa(maxID))
	}

	utils.RespondJSON(w, http.StatusOK, result)
}