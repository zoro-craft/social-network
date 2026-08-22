package posts

import (
	"net/http"
	"strconv"

	"zone/internal/auth"
	"zone/internal/database"
	"zone/internal/utils"
)

const postSelect = `
	SELECT p.id, p.title, p.content, COALESCE(p.image, ''), p.privacy,
	       p.user_id, COALESCE(u.nickname, ''), p.created_at
	FROM posts p
	JOIN users u ON p.user_id = u.id
`

func GetPostsAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	viewerID, _ := auth.GetUserIDFromSession(r)

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 50 {
		limit = 10
	}
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if offset < 0 {
		offset = 0
	}

	query := postSelect + "WHERE " + visibleCondition + " ORDER BY p.id DESC LIMIT ? OFFSET ?"
	rows, err := database.Database.Query(query, viewerID, viewerID, viewerID, limit, offset)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Database error loading posts")
		return
	}
	defer rows.Close()

	result := []Post{}
	for rows.Next() {
		var p Post
		if err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.Image, &p.Privacy, &p.UserID, &p.Nickname, &p.CreatedAt); err == nil {
			result = append(result, p)
		}
	}

	utils.RespondJSON(w, http.StatusOK, result)
}

func GetPostsByUser(userID, viewerID int) ([]Post, error) {
	query := postSelect + "WHERE p.user_id = ? AND " + visibleCondition + " ORDER BY p.id DESC"
	rows, err := database.Database.Query(query, userID, viewerID, viewerID, viewerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []Post{}
	for rows.Next() {
		var p Post
		if err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.Image, &p.Privacy, &p.UserID, &p.Nickname, &p.CreatedAt); err == nil {
			result = append(result, p)
		}
	}
	return result, rows.Err()
}