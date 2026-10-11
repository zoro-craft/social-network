package reels

import (
	"net/http"

	"zone/internal/auth"
	"zone/internal/database"
	"zone/internal/utils"
)

func DeleteReelAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
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

	// Scoping the delete to the author prevents removing someone else's reel.
	// Likes, comments, views, saves and shares cascade away with it.
	result, err := database.Database.Exec(
		"DELETE FROM reels WHERE id = ? AND user_id = ?",
		reelID, userID,
	)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to delete reel")
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		utils.HandleError(w, http.StatusNotFound, "Reel not found")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Reel deleted"})
}
