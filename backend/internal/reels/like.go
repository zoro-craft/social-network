package reels

import (
	"net/http"
	"time"

	"zone/internal/auth"
	"zone/internal/database"
	"zone/internal/utils"
)

// toggleReelRow flips one (reel_id, user_id) row in a single-row table: an
// existing row is deleted, a missing one is inserted. `table` is always a
// compile-time constant from this package, never user input.
func toggleReelRow(table string, reelID, userID int) error {
	result, err := database.Database.Exec(
		"DELETE FROM "+table+" WHERE reel_id = ? AND user_id = ?",
		reelID, userID,
	)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected > 0 {
		return nil
	}

	_, err = database.Database.Exec(
		"INSERT INTO "+table+" (reel_id, user_id, created_at) VALUES (?, ?, ?)",
		reelID, userID, time.Now().UnixMilli(),
	)
	return err
}

// ToggleReelLikeAPI likes an unliked reel and unlikes a liked one, mirroring the
// click-again-to-undo behaviour of post reactions.
func ToggleReelLikeAPI(w http.ResponseWriter, r *http.Request) {
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

	if err := toggleReelRow("reel_likes", reelID, userID); err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Could not save like")
		return
	}

	respondWithStats(w, reelID, userID)
}
