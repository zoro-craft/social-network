package stories

import (
	"net/http"
	"strconv"

	"zone/internal/auth"
	"zone/internal/database"
	"zone/internal/utils"
)

func DeleteStoryAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID, err := auth.GetUserID(r)
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, "You must be logged in")
		return
	}

	storyID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || storyID <= 0 {
		utils.HandleError(w, http.StatusBadRequest, "Invalid story id")
		return
	}

	// Scoping the delete to the author prevents removing someone else's story.
	result, err := database.Database.Exec(
		"DELETE FROM stories WHERE id = ? AND user_id = ?",
		storyID, userID,
	)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to delete story")
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		utils.HandleError(w, http.StatusNotFound, "Story not found")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Story deleted"})
}
