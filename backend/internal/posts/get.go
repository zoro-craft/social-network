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

	viewerID, err := auth.GetUserID(r)
	if err!=nil{
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
	if beforeStr == "" {
		var maxID int
		if err := database.Database.QueryRow("SELECT COALESCE(MAX(id), 0) FROM posts").Scan(&maxID); err == nil {
			w.Header().Set("X-Max-Post-Id", strconv.Itoa(maxID))
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