package handlers

import (
	"net/http"
	"strconv"

	"zone/internal/auth"
	"zone/internal/comments/models"
	"zone/internal/database"
	"zone/internal/posts/repositories"
	"zone/internal/utils"
)

const commentSelect = `
	SELECT c.id, c.post_id, c.group_id, c.content, COALESCE(c.image, ''), c.user_id,
	       COALESCE(u.first_name, ''), COALESCE(u.avatar, ''), c.created_at
	FROM comments c
	JOIN users u ON c.user_id = u.id
`

func GetCommentsAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	postID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || postID <= 0 {
		utils.HandleError(w, http.StatusBadRequest, "Invalid post id")
		return
	}

	viewerID, err := auth.GetUserID(r)
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, "You must be logged in")
		return
	}

	if !repositories.CanViewPost(postID, viewerID) {
		utils.HandleError(w, http.StatusNotFound, "Post not found")
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
		WHERE c.post_id = ?
	`
	args := []interface{}{postID}

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

	result := []models.Comment{}
	for rows.Next() {
		var c models.Comment
		if err := rows.Scan(&c.ID, &c.PostID, &c.GroupID, &c.Content, &c.Image, &c.UserID, &c.FirstName, &c.Avatar, &c.CreatedAt); err == nil {
			result = append(result, c)
		}
	}

	if beforeStr == "" {
		var maxID int
		database.Database.QueryRow(
			"SELECT COALESCE(MAX(id), 0) FROM comments WHERE post_id = ?", postID,
		).Scan(&maxID)
		w.Header().Set("X-Max-Comment-Id", strconv.Itoa(maxID))
	}

	utils.RespondJSON(w, http.StatusOK, result)
}
