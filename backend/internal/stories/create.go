package stories

import (
	"net/http"
	"strings"
	"time"

	"zone/internal/auth"
	"zone/internal/database"
	"zone/internal/utils"
)

const maxCaptionLength = 200

func CreateStoryAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID, err := auth.GetUserID(r)
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, "You must be logged in to post a story")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, utils.MaxImageFormSize)
	if err := r.ParseMultipartForm(utils.MaxUploadSize); err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Invalid form data or file too large")
		return
	}

	caption := strings.TrimSpace(r.FormValue("caption"))
	if len(caption) > maxCaptionLength {
		utils.HandleError(w, http.StatusBadRequest, "Caption too long")
		return
	}

	image, err := utils.SaveImage(r, "stories", "/uploads/stories/")
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, err.Error())
		return
	}

	// A story is media-first: with no image there is nothing to show.
	if image == "" {
		utils.HandleError(w, http.StatusBadRequest, "An image is required for a story")
		return
	}

	_, err = database.Database.Exec(
		"INSERT INTO stories (user_id, image, caption, created_at) VALUES (?, ?, ?, ?)",
		userID, image, caption, time.Now().UnixMilli(),
	)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to create story")
		return
	}

	utils.RespondJSON(w, http.StatusCreated, map[string]string{"message": "Story created successfully!"})
}
