package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"zone/internal/auth"
	"zone/internal/comments/models"
	"zone/internal/database"
	"zone/internal/posts/repositories"
	"zone/internal/utils"
)

func CreateCommentAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID, err := auth.GetUserID(r)
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, "You must be logged in to comment")
		return
	}

	postID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || postID <= 0 {
		utils.HandleError(w, http.StatusBadRequest, "Invalid post id")
		return
	}

	if !repositories.CanViewPost(postID, userID) {
		utils.HandleError(w, http.StatusBadRequest, "Invalid post")
		return
	}
	var groupID *int
	err = database.Database.QueryRow(
		"SELECT group_id FROM posts WHERE id = ?", postID,
	).Scan(&groupID)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to load post")
		return
	}

	if groupID != nil {
		var isMember bool
		err = database.Database.QueryRow(
			"SELECT EXISTS(SELECT * FROM group_members WHERE group_id = ? AND user_id = ?)",
			*groupID, userID,
		).Scan(&isMember)
		if err != nil || !isMember {
			utils.HandleError(w, http.StatusForbidden, "You must be a group member to comment")
			return
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, utils.MaxImageFormSize)
	if err := r.ParseMultipartForm(utils.MaxUploadSize); err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Invalid form data or file too large")
		return
	}

	c := models.Comment{
		PostID:  postID,
		GroupID: groupID,
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

	c.Image, err = utils.SaveImage(r, "comments", "/uploads/comments/")
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := database.Database.Exec(
		"INSERT INTO comments (post_id, group_id, user_id, content, image) VALUES (?, ?, ?, ?, ?)",
		c.PostID, c.GroupID, c.UserID, c.Content, c.Image,
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
		SELECT c.id, c.post_id, c.group_id, c.content, COALESCE(c.image, ''), c.user_id,
		u.first_name, COALESCE(u.avatar, ''), c.created_at
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.id = ?
	`, commentID).Scan(&c.ID, &c.PostID, &c.GroupID, &c.Content, &c.Image, &c.UserID, &c.FirstName, &c.Avatar, &c.CreatedAt)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to load created comment")
		return
	}

	utils.RespondJSON(w, http.StatusCreated, c)
}
