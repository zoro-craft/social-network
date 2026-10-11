package reels

import (
	"net/http"
	"time"

	"zone/internal/auth"
	"zone/internal/database"
	"zone/internal/utils"
)

// ShareReelAPI records a share. Unlike likes and saves this never toggles: every
// call adds to the count, since sharing the same reel twice is still two shares.
func ShareReelAPI(w http.ResponseWriter, r *http.Request) {
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
		"INSERT INTO reel_shares (reel_id, user_id, created_at) VALUES (?, ?, ?)",
		reelID, userID, time.Now().UnixMilli(),
	)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Could not record share")
		return
	}

	respondWithStats(w, reelID, userID)
}
