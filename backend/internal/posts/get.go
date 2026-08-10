package posts

import (
	"net/http"
	"strconv"

	"zone/internal/auth"
	"zone/internal/database"
	"zone/internal/utils"
)

func GetPostsAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	viewerID, _ := auth.GetUserIDFromSession(r)

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil || offset < 0 {
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

	// A post is visible if it's public, or if it's "almost_private"/"private"
	// and the viewer is logged in. (Follower-only filtering comes later.)
	query := `
		SELECT p.id, p.title, p.content, COALESCE(p.image, ''), p.privacy,
		       p.user_id, COALESCE(u.nickname, ''), p.created_at
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE (p.privacy = 'public' OR ? > 0)
	`
	args := []interface{}{viewerID}

	if beforeStr != "" {
		query += " AND p.id <= ?"
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

	posts := []Post{}
	for rows.Next() {
		var p Post
		if err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.Image, &p.Privacy, &p.UserID, &p.Nickname, &p.CreatedAt); err != nil {
			continue
		}
		posts = append(posts, p)
	}

	if beforeStr == "" {
		var maxID int
		if err := database.Database.QueryRow("SELECT COALESCE(MAX(id), 0) FROM posts").Scan(&maxID); err == nil {
			w.Header().Set("X-Max-Post-Id", strconv.Itoa(maxID))
		}
	}

	utils.RespondJSON(w, http.StatusOK, posts)
}


func GetPostsByUser(userID, viewerID int) ([]Post, error) {
	rows, err := database.Database.Query(`
		SELECT p.id, p.title, p.content, COALESCE(p.image, ''), p.privacy,
		       p.user_id, COALESCE(u.nickname, ''), p.created_at
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.user_id = ? AND (p.privacy = 'public' OR ? > 0)
		ORDER BY p.id DESC
	`, userID, viewerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []Post{}
	for rows.Next() {
		var p Post
		if err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.Image, &p.Privacy, &p.UserID, &p.Nickname, &p.CreatedAt); err != nil {
			continue
		}
		result = append(result, p)
	}
	return result, rows.Err()
}