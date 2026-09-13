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

	userID, err := auth.GetUserID(r)
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, "You must be logged in to post")
		return
	}
	
	r.Body = http.MaxBytesReader(w, r.Body, utils.MaxUploadSize)
	if err := r.ParseMultipartForm(utils.MaxUploadSize); err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Invalid form data or file too large")
		return
	}

	p := Post{
		Title:   strings.TrimSpace(r.FormValue("title")),
		Content: strings.TrimSpace(r.FormValue("content")),
		Privacy: strings.TrimSpace(r.FormValue("privacy")),
		UserID:  userID,
	}

//group
	groupIDStr := strings.TrimSpace(r.FormValue("group_id"))

	if groupIDStr != "" {
		groupID, err := strconv.Atoi(groupIDStr)
		if err != nil || groupID <= 0 {
			utils.HandleError(w, http.StatusBadRequest, "Invalid group id")
			return
		}

		// check goup
		var exists bool

		err = database.Database.QueryRow(
			"SELECT EXISTS(SELECT 1 FROM user_groups WHERE id = ?)",
			groupID,
		).Scan(&exists)

		if err != nil || !exists {
			utils.HandleError(w, http.StatusBadRequest, "Group not found")
			return
		}
		if !utils.IsGroupMember(groupID, userID) {
			utils.HandleError(w, http.StatusForbidden, "You must be a group member to post")
			return
		}

		p.GroupID = &groupID

		// group posts are always public
		p.Privacy = "public"
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

	if p.GroupID == nil {
		if p.Privacy == "" {
			p.Privacy = "public"
		}
		if !allowedPrivacy[p.Privacy] {
			utils.HandleError(w, http.StatusBadRequest, "Invalid privacy option")
			return
		}
	}

	viewerIDs := r.Form["viewer_ids"]

	if p.GroupID == nil &&
		p.Privacy == "private" &&
		len(viewerIDs) == 0 {
		utils.HandleError(w, http.StatusBadRequest, "Select at least one follower")
		return
	}

	if p.GroupID == nil {
		for _, idStr := range viewerIDs {
			viewerID, err := strconv.Atoi(idStr)
			if err != nil {
				utils.HandleError(w, http.StatusBadRequest, "Invalid viewer id")
				return
			}

			var isFollower bool
			err = database.Database.QueryRow(
				"SELECT EXISTS(SELECT 1 FROM followers WHERE follower_id = ? AND followee_id = ? AND status = 'accepted')",
				viewerID, userID,
			).Scan(&isFollower)

			if err != nil {
				utils.HandleError(w, http.StatusInternalServerError, "Database error")
				return
			}

			if !isFollower {
				utils.HandleError(w, http.StatusBadRequest, "Viewer must be an accepted follower")
				return
			}
		}
	}

	p.Image, err = utils.SaveImage(r, "./uploads/posts", "/uploads/posts/")

	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := database.Database.Exec(`INSERT INTO posts (user_id, group_id, title, content, image, privacy) VALUES (?, ?, ?, ?, ?, ?)`,
		p.UserID, p.GroupID, p.Title, p.Content, p.Image, p.Privacy,
	)

	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to create post")
		return
	}

	postID, _ := result.LastInsertId()

	if p.GroupID == nil {
		for _, idStr := range viewerIDs {
			viewerID, _ := strconv.Atoi(idStr)

			_, err = database.Database.Exec(
				"INSERT OR IGNORE INTO post_viewers (post_id, user_id) VALUES (?, ?)",
				postID, viewerID,
			)
			if err != nil {
				utils.HandleError(w, http.StatusInternalServerError, "Failed to save post viewers")
				return
			}
		}
	}

	utils.RespondJSON(w, http.StatusCreated, map[string]string{"message": "Post created successfully!"})

}