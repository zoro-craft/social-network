package comments

import (
	"net/http"
	"strconv"
	"strings"

	"zone/internal/auth"
	"zone/internal/database"
	"zone/internal/utils"
)

func CreateCommentAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID, err := auth.GetUserIDFromSession(r)
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, "You must be logged in to comment")
		return
	}

	postID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || postID <= 0 {
		utils.HandleError(w, http.StatusBadRequest, "Invalid post id")
		return
	}

	var exists bool
	if err := database.Database.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM posts WHERE id = ?)", postID,
	).Scan(&exists); err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Database error")
		return
	}
	if !exists {
		utils.HandleError(w, http.StatusBadRequest, "Invalid post")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, utils.MaxCommentUploadSize)
	if err := r.ParseMultipartForm(utils.MaxCommentUploadSize); err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Invalid form data or file too large")
		return
	}

	c := Comment{
		PostID:  postID,
		UserID:  userID,
		Content: strings.TrimSpace(r.FormValue("content")),
	}

	if c.Content == "" {
		utils.HandleError(w, http.StatusBadRequest, "Content is required")
		return
	}
	if len(c.Content) > 1000 {
		utils.HandleError(w, http.StatusBadRequest, "Content too long")
		return
	}

	c.Image, err = saveCommentImage(r)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := database.Database.Exec(
		"INSERT INTO comments (post_id, user_id, content, image) VALUES (?, ?, ?, ?)",
		c.PostID, c.UserID, c.Content, c.Image,
	)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to create comment")
		return
	}

	commentID, err := result.LastInsertId()
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Could not get new comment ID")
		return
	}

	err = database.Database.QueryRow(`
		SELECT c.id, c.post_id, c.content, COALESCE(c.image, ''), c.user_id,
		COALESCE(u.first_name, ''), COALESCE(u.avatar, ''), c.created_at
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.id = ?
	`, commentID).Scan(&c.ID, &c.PostID, &c.Content, &c.Image, &c.UserID, &c.FirstName, &c.Avatar, &c.CreatedAt)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to load created comment")
		return
	}

	utils.RespondJSON(w, http.StatusCreated, c)
}