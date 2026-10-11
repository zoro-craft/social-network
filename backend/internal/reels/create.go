package reels

import (
	"net/http"
	"strings"
	"time"

	"zone/internal/auth"
	"zone/internal/database"
	"zone/internal/utils"
)

const maxCaptionLength = 300

func CreateReelAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID, err := auth.GetUserID(r)
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, "You must be logged in to post a reel")
		return
	}

	// Video bodies are large, so cap the stream and let ParseMultipartForm spill
	// anything past MultipartMemory to a temp file rather than the heap.
	r.Body = http.MaxBytesReader(w, r.Body, utils.MaxVideoUploadSize)
	if err := r.ParseMultipartForm(utils.MultipartMemory); err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Invalid form data or file too large")
		return
	}

	caption := strings.TrimSpace(r.FormValue("caption"))
	if len(caption) > maxCaptionLength {
		utils.HandleError(w, http.StatusBadRequest, "Caption too long")
		return
	}

	video, err := utils.SaveVideo(r, "video", "reels", "/uploads/reels/")
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, err.Error())
		return
	}

	// A reel is media-first: with no video there is nothing to play.
	if video == "" {
		utils.HandleError(w, http.StatusBadRequest, "A video is required for a reel")
		return
	}

	result, err := database.Database.Exec(
		"INSERT INTO reels (user_id, video, caption, created_at) VALUES (?, ?, ?, ?)",
		userID, video, caption, time.Now().UnixMilli(),
	)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to create reel")
		return
	}

	reelID, err := result.LastInsertId()
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to create reel")
		return
	}

	reel, err := GetReel(int(reelID), userID)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to load the new reel")
		return
	}

	utils.RespondJSON(w, http.StatusCreated, reel)
}
