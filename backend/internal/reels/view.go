package reels

import (
	"net/http"
	"time"

	"zone/internal/auth"
	"zone/internal/database"
	"zone/internal/utils"
)

// ReelViewAPI counts a play. Replays must not inflate the number, so the first
// view from a given user wins — the same rule stories use for seen markers.
func ReelViewAPI(w http.ResponseWriter, r *http.Request) {
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

	_, err = database.Database.Exec(
		"INSERT OR IGNORE INTO reel_views (reel_id, user_id, viewed_at) VALUES (?, ?, ?)",
		reelID, userID, time.Now().UnixMilli(),
	)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Could not record view")
		return
	}

	respondWithStats(w, reelID, userID)
}
