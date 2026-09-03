package profile

import (
	"net/http"

	"zone/internal/auth"
	"zone/internal/database"
	"zone/internal/utils"
)

func ToggleVisibilityAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID, err := auth.GetUserID(r)
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	_, err = database.Database.Exec(
		"UPDATE users SET is_private = NOT is_private WHERE id = ?", userID,
	)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to update profile")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Profile visibility updated"})
}
