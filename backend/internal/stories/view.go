package stories

import (
	"net/http"
	"strconv"
	"time"

	"zone/internal/auth"
	"zone/internal/database"
	"zone/internal/utils"
)

func StoryViewAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	if !CanViewStory(storyID, userID) {
		utils.HandleError(w, http.StatusNotFound, "Story not found")
		return
	}

	// Re-watching a story must not fail, so the first view wins.
	_, err = database.Database.Exec(
		"INSERT OR IGNORE INTO story_views (story_id, user_id, viewed_at) VALUES (?, ?, ?)",
		storyID, userID, time.Now().UnixMilli(),
	)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to mark story as seen")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Story marked as seen"})
}
