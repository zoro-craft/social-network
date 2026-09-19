package posts

import (
	"net/http"
	"strconv"

	"zone/internal/auth"
	"zone/internal/database"
	"zone/internal/utils"
)

const postSelect = `
	SELECT
		p.id,
		p.title,
		p.content,
		COALESCE(p.image, ''),
		p.privacy,
		p.group_id,
		p.user_id,
		u.full_name,
		p.created_at
`

const reactionSelect = `,
	(SELECT COUNT(*)
	 FROM post_reactions
	 WHERE post_id = p.id AND is_like = 1) AS like_count,

	(SELECT COUNT(*)
	 FROM post_reactions
	 WHERE post_id = p.id AND is_like = 0) AS dislike_count
`

const fromPosts = `
	FROM posts p
	JOIN users u ON p.user_id = u.id
`

func GetPostsAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	viewerID, err := auth.GetUserID(r)
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, "You must be logged in")
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

	// Home posts + reactions
	query := postSelect + reactionSelect + fromPosts + " WHERE p.group_id IS NULL AND (" + visibleCondition + ")"

	args := []interface{}{viewerID, viewerID, viewerID, viewerID}

	if beforeID > 0 {
		query += " AND p.id < ?"
		args = append(args, beforeID)
	}

	query += " ORDER BY p.id DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := database.Database.Query(query, args...)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Database error loading posts")
		return
	}
	defer rows.Close()

	result := []Post{}

	for rows.Next() {
		var p Post

		if err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.Image, &p.Privacy, &p.GroupID, &p.UserID, &p.FullName, &p.CreatedAt, &p.LikeCount, &p.DislikeCount); err != nil {
			utils.HandleError(w, http.StatusInternalServerError, "Error scanning post")
			return
		}

		result = append(result, p)
	}

	if err := rows.Err(); err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Error reading posts")
		return
	}

	if beforeStr == "" {
		var maxID int

		if err := database.Database.QueryRow("SELECT COALESCE(MAX(id), 0) FROM posts").Scan(&maxID); err == nil {
			w.Header().Set("X-Max-Post-Id", strconv.Itoa(maxID))
		}
	}

	utils.RespondJSON(w, http.StatusOK, result)
}

func GetPostsByUser(userID, viewerID int) ([]Post, error) {
	query := postSelect + fromPosts + " WHERE p.user_id = ? AND " + visibleCondition + " ORDER BY p.id DESC"

	rows, err := database.Database.Query(query, userID, viewerID, viewerID, viewerID, viewerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []Post{}

	for rows.Next() {
		var p Post

		if err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.Image, &p.Privacy, &p.GroupID, &p.UserID, &p.FullName, &p.CreatedAt); err != nil {
			return nil, err
		}

		result = append(result, p)
	}

	return result, rows.Err()
}

func GetPostsByGroup(groupID int) ([]Post, error) {
	query := postSelect + reactionSelect + fromPosts + " WHERE p.group_id = ? ORDER BY p.id DESC"

	rows, err := database.Database.Query(query, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []Post{}

	for rows.Next() {
		var p Post

		if err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.Image, &p.Privacy, &p.GroupID, &p.UserID, &p.FullName, &p.CreatedAt, &p.LikeCount, &p.DislikeCount); err != nil {
			return nil, err
		}

		result = append(result, p)
	}

	return result, rows.Err()
}

func GetPostsByGroupAPI(w http.ResponseWriter, r *http.Request) {
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

	posts, err := GetPostsByGroup(groupID)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to load group posts")
		return
	}

	utils.RespondJSON(w, http.StatusOK, posts)
}
