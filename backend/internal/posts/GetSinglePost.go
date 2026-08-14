package posts

import (
	"net/http"
	"strconv"
	"zone/internal/database"
	"zone/internal/utils"
)

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
	err = database.Database.QueryRow(`
		SELECT p.id, p.title, p.content, p.image, p.privacy, p.user_id, u.nickname, p.created_at
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.id = ?
	`, postID).Scan(&p.ID, &p.Title, &p.Content, &p.Image, &p.Privacy, &p.UserID, &p.Nickname, &p.CreatedAt)
	if err != nil {
		utils.HandleError(w, http.StatusNotFound, "Post not found")
		return
	}

	utils.RespondJSON(w, http.StatusOK, p)
}