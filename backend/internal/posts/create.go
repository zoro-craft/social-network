package posts

import (
	"errors"
	"net/http"
	"strings"

	"zone/internal/auth"
	"zone/internal/database"
	"zone/internal/utils"
)

const maxPostUploadSize = 5 << 20 // 5MB

var allowedPostExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
}

var allowedPrivacy = map[string]bool{
	"public":         true,
	"almost_private": true,
	"private":        true,
}

var errInvalidImageType = errors.New("only JPEG, PNG, and GIF images are allowed")

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

	r.Body = http.MaxBytesReader(w, r.Body, maxPostUploadSize)
	if err := r.ParseMultipartForm(maxPostUploadSize); err != nil {
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
	// NOTE: almost_private and private are accepted and stored so the
	// frontend can offer all three choices, but visibility filtering
	// for them isn't implemented yet — every post is currently visible
	// to everyone regardless of this value.

	p.Image, err = savePostImage(r)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, err.Error())
		return
	}

_, err = database.Database.Exec(
	"INSERT INTO posts (user_id, title, content, image, privacy) VALUES (?, ?, ?, ?, ?)",
	p.UserID, p.Title, p.Content, p.Image, p.Privacy,
)

	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to create post")
		return
	}

	utils.RespondJSON(w, http.StatusCreated, map[string]string{"message": "Post created successfully!"})
}
