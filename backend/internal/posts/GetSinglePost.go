package posts

import (
	"net/http"
	"strconv"
	"zone/internal/auth"
	"zone/internal/database"
	"zone/internal/utils"
)

func GetPostByIDAPI(w http.ResponseWriter, r *http.Request) {
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
	var p Post
	
	query := `
		SELECT p.id, p.title, p.content, p.image, p.privacy, p.user_id, u.nickname, p.created_at
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.id = ? AND ` + visibleCondition

	err = database.Database.QueryRow(query, postID, viewerID, viewerID, viewerID).Scan(
		&p.ID, &p.Title, &p.Content, &p.Image, &p.Privacy, &p.UserID, &p.Nickname, &p.CreatedAt,
	)
	
	if err != nil {
		utils.HandleError(w, http.StatusNotFound, "Post not found or unauthorized")
		return
	}

	utils.RespondJSON(w, http.StatusOK, p)
}