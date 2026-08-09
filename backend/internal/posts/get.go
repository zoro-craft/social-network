package posts

import (
	"net/http"
	"strconv"

	"zone/internal/database"
	"zone/internal/utils"
)

func GetPostsAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

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

	query := `
SELECT p.id, p.title, p.content, p.image, p.privacy, p.user_id, u.nickname, p.created_at	FROM posts p
		JOIN users u ON p.user_id = u.id
	`
	args := []interface{}{}

	if beforeStr != "" {
		query += " WHERE p.id <= ?"
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
		var image *string

		if err :=rows.Scan(&p.ID, &p.Title, &p.Content, &image, &p.Privacy, &p.UserID, &p.Nickname, &p.CreatedAt); err != nil {
			continue
		}
		if image != nil {
			p.Image = *image
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
/* 
func GetPostByIDAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	idStr := r.PathValue("id")
	postID, err := strconv.Atoi(idStr)
	if err != nil || postID <= 0 {
		utils.HandleError(w, http.StatusBadRequest, "Invalid post id")
		return
	}

	var p Post
	var image *string

	err = database.Database.QueryRow(`
		SELECT p.id, p.content, p.image, p.privacy, p.user_id, u.nickname, p.created_at
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.id = ?
	`, postID).Scan(&p.ID, &p.Content, &image, &p.Privacy, &p.UserID, &p.Nickname, &p.CreatedAt)
	if err != nil {
		utils.HandleError(w, http.StatusNotFound, "Post not found")
		return
	}
	if image != nil {
		p.Image = *image
	}

	utils.RespondJSON(w, http.StatusOK, p)
} */