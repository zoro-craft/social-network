package comments

import (
	"net/http"
	"strconv"

	"zone/internal/auth"
	"zone/internal/database"
	"zone/internal/posts"
	"zone/internal/utils"
)

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

	viewerID, _ := auth.GetUserIDFromSession(r)

	if !posts.CanViewPost(postID, viewerID) {
		utils.HandleError(w, http.StatusNotFound, "Post not found")
		return
	}

	rows, err := database.Database.Query(`
		SELECT c.id, c.post_id, c.content, COALESCE(c.image, ''), c.user_id,
		       COALESCE(u.first_name, ''), COALESCE(u.avatar, ''), c.created_at
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.post_id = ?
		ORDER BY c.id DESC
	`, postID)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Database error loading comments")
		return
	}
	defer rows.Close()

	result := []Comment{}
	for rows.Next() {
		var c Comment
		if rows.Scan(&c.ID, &c.PostID, &c.Content, &c.Image, &c.UserID, &c.FirstName, &c.Avatar, &c.CreatedAt) == nil {
		continue
		}
		result = append(result, c)
	}

	utils.RespondJSON(w, http.StatusOK, result)
}