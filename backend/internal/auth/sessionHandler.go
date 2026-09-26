package auth

import (
	"fmt"
	"net/http"
	"zone/internal/utils"

	"zone/internal/auth/repositories"
	"zone/internal/middleware"
)

func GetUserID(r *http.Request) (int, error) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		return 0, fmt.Errorf("user ID not found in context")
	}
	return userID, nil
}

func MeHandler(w http.ResponseWriter, r *http.Request) {

	userID, err := GetUserID(r)
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, "Unauthorized. Please login.")
		return
	}

	profile, err := repositories.ProfileByID(userID)
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, "User not found or database error.")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"status": "success",
		"user": map[string]interface{}{
			"id":         userID,
			"nickname":   profile.Nickname,
			"email":      profile.Email,
			"first_name": profile.FirstName,
			"last_name":  profile.LastName,
			"avatar":     profile.Avatar,
		},
	})
}
