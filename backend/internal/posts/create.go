package posts

import (
	"net/http"
	"strconv"
	"strings"

	"zone/internal/auth"
	"zone/internal/database"
	"zone/internal/utils"
)

var allowedPrivacy = map[string]bool{
	"public":         true,
	"almost_private": true,
	"private":        true,
}

func CreatePostAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID, err := auth.GetUserIDFromSession(r)
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, "You must be logged in to post")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, utils.MaxCommentUploadSize)
	if err := r.ParseMultipartForm(utils.MaxCommentUploadSize); err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Invalid form data or file too large")
		return
	}

	p := Post{
		Title:   strings.TrimSpace(r.FormValue("title")),
		Content: strings.TrimSpace(r.FormValue("content")),
		Privacy: strings.TrimSpace(r.FormValue("privacy")),
		UserID:  userID,
	}

	if p.Title == "" {
		utils.HandleError(w, http.StatusBadRequest, "Title is required")
		return
	}
		if len(p.Title) > 100 {
		utils.HandleError(w, http.StatusBadRequest, "Title too long")
		return
	}

	if p.Content == "" {
		utils.HandleError(w, http.StatusBadRequest, "Content is required")
		return
	}
	if len(p.Content) > 2000 {
		utils.HandleError(w, http.StatusBadRequest, "Content too long")
		return
	}
	if p.Privacy == "" {
		p.Privacy = "public"
	}
	if !allowedPrivacy[p.Privacy] {
		utils.HandleError(w, http.StatusBadRequest, "Invalid privacy option")
		return
	}

	viewerIDs := r.Form["viewer_ids"]
	if p.Privacy == "private" && len(viewerIDs) == 0 {
		utils.HandleError(w, http.StatusBadRequest, "Select at least one follower")
		return
	}

	p.Image, err = savePostImage(r)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := database.Database.Exec(
		"INSERT INTO posts (user_id, title, content, image, privacy) VALUES (?, ?, ?, ?, ?)",
		p.UserID, p.Title, p.Content, p.Image, p.Privacy,
	)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to create post")
		return
	}

	postID, _ := result.LastInsertId()

	for _, idStr := range viewerIDs {
		viewerID, err := strconv.Atoi(idStr)
		if err == nil {
			database.Database.Exec(
				"INSERT OR IGNORE INTO post_viewers (post_id, user_id) VALUES (?, ?)",
				postID, viewerID,
			)
		}
	}

	utils.RespondJSON(w, http.StatusCreated, map[string]string{"message": "Post created successfully!"})
}

