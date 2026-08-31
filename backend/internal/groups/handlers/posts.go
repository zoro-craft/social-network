package handlers

import (
	"errors"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"zone/internal/auth"
	"zone/internal/groups/repositories"
	"zone/internal/groups/services"
	"zone/internal/utils"
)

func saveGroupPostImage(r *http.Request) (string, error) {
	file, header, err := r.FormFile("image")
	if err == http.ErrMissingFile {
		return "", nil
	}
	if err != nil {
		return "", nil
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !utils.AllowedCommentExtensions[ext] {
		return "", utils.ErrInvalidImageType
	}

	filename, err := utils.SaveUploadFile(file, header.Filename, "./uploads/posts")
	if err != nil {
		return "", err
	}
	return "/uploads/posts/" + filename, nil
}

func requireGroupMember(w http.ResponseWriter, r *http.Request) (groupID, userID int, ok bool) {
	userID, err := auth.GetUserIDFromSession(r)
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, "You must be logged in")
		return 0, 0, false
	}

	groupID, err = strconv.Atoi(r.PathValue("id"))
	if err != nil || groupID <= 0 {
		utils.HandleError(w, http.StatusBadRequest, "Invalid group id")
		return 0, 0, false
	}

	if !services.GroupExists(groupID) {
		utils.HandleError(w, http.StatusNotFound, "Group not found")
		return 0, 0, false
	}

	if !services.IsGroupMember(groupID, userID) && !services.IsGroupCreator(groupID, userID) {
		utils.HandleError(w, http.StatusForbidden, "You must be a group member")
		return 0, 0, false
	}

	return groupID, userID, true
}

// CreateGroupPostHandler lets a member create a post inside a group.
func CreateGroupPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	groupID, userID, ok := requireGroupMember(w, r)
	if !ok {
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, utils.MaxCommentUploadSize)
	if err := r.ParseMultipartForm(utils.MaxCommentUploadSize); err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Invalid form data or file too large")
		return
	}

	title := strings.TrimSpace(r.FormValue("title"))
	content := strings.TrimSpace(r.FormValue("content"))

	if title == "" {
		utils.HandleError(w, http.StatusBadRequest, "Title is required")
		return
	}
	if content == "" {
		utils.HandleError(w, http.StatusBadRequest, "Content is required")
		return
	}
	if len(title) > 100 {
		utils.HandleError(w, http.StatusBadRequest, "Title too long")
		return
	}
	if len(content) > 5000 {
		utils.HandleError(w, http.StatusBadRequest, "Content too long")
		return
	}

	image, err := saveGroupPostImage(r)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, err.Error())
		return
	}

	postID, err := repositories.CreateGroupPost(groupID, userID, title, content, image)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to create group post")
		return
	}

	post, err := repositories.GetGroupPost(int(postID))
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to load created post")
		return
	}

	utils.RespondJSON(w, http.StatusCreated, post)
}

// GroupPostsHandler returns all posts in a group (members only).
func GroupPostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	groupID, _, ok := requireGroupMember(w, r)
	if !ok {
		return
	}

	posts, err := repositories.ListGroupPosts(groupID)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to load group posts")
		return
	}

	utils.RespondJSON(w, http.StatusOK, posts)
}

// CreateGroupCommentHandler lets a member comment on a group post.
func CreateGroupCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID, err := auth.GetUserIDFromSession(r)
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, "You must be logged in")
		return
	}

	postID, err := strconv.Atoi(r.PathValue("postID"))
	if err != nil || postID <= 0 {
		utils.HandleError(w, http.StatusBadRequest, "Invalid post id")
		return
	}

	groupID, err := repositories.GroupIDForPost(postID)
	if err != nil || groupID <= 0 {
		utils.HandleError(w, http.StatusNotFound, "Group post not found")
		return
	}

	if !services.IsGroupMember(groupID, userID) && !services.IsGroupCreator(groupID, userID) {
		utils.HandleError(w, http.StatusForbidden, "You must be a group member")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, utils.MaxCommentUploadSize)
	if err := r.ParseMultipartForm(utils.MaxCommentUploadSize); err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Invalid form data or file too large")
		return
	}

	content := strings.TrimSpace(r.FormValue("content"))
	if content == "" {
		utils.HandleError(w, http.StatusBadRequest, "Content is required")
		return
	}
	if len(content) > 1000 {
		utils.HandleError(w, http.StatusBadRequest, "Content too long")
		return
	}

	image, err := saveGroupCommentImage(r)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, err.Error())
		return
	}

	commentID, err := repositories.CreateGroupComment(postID, userID, groupID, content, image)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to create comment")
		return
	}

	utils.RespondJSON(w, http.StatusCreated, map[string]int{"id": int(commentID), "post_id": postID})
}

// GroupCommentsHandler returns comments on a group post (members only).
func GroupCommentsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID, err := auth.GetUserIDFromSession(r)
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, "You must be logged in")
		return
	}

	postID, err := strconv.Atoi(r.PathValue("postID"))
	if err != nil || postID <= 0 {
		utils.HandleError(w, http.StatusBadRequest, "Invalid post id")
		return
	}

	groupID, err := repositories.GroupIDForPost(postID)
	if err != nil || groupID <= 0 {
		utils.HandleError(w, http.StatusNotFound, "Group post not found")
		return
	}

	if !services.IsGroupMember(groupID, userID) && !services.IsGroupCreator(groupID, userID) {
		utils.HandleError(w, http.StatusForbidden, "You must be a group member")
		return
	}

	comments, err := repositories.ListGroupComments(postID)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to load comments")
		return
	}

	utils.RespondJSON(w, http.StatusOK, comments)
}

func saveGroupCommentImage(r *http.Request) (string, error) {
	file, header, err := r.FormFile("image")
	if err != nil {
		return "", nil
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !utils.AllowedCommentExtensions[ext] {
		return "", errors.New(utils.ErrInvalidImageType.Error())
	}

	filename, err := utils.SaveUploadFile(file, header.Filename, "./uploads/comments")
	if err != nil {
		return "", err
	}
	return "/uploads/comments/" + filename, nil
}
