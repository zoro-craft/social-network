package reels

import (
	"net/http"

	"zone/internal/auth"
	"zone/internal/utils"
)

// ToggleReelSaveAPI bookmarks a reel for the viewer, or removes the bookmark.
func ToggleReelSaveAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID, err := auth.GetUserID(r)
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, "You must be logged in")
		return
	}

	reelID := reelIDFromPath(r)
	if reelID == 0 {
		utils.HandleError(w, http.StatusBadRequest, "Invalid reel id")
		return
	}

	if !reelExists(reelID) {
		utils.HandleError(w, http.StatusNotFound, "Reel not found")
		return
	}

	if err := toggleReelRow("reel_saves", reelID, userID); err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Could not save reel")
		return
	}

	respondWithStats(w, reelID, userID)
}
