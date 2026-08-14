package followers

import (
	"net/http"

	"zone/internal/auth"
	"zone/internal/database"
	"zone/internal/utils"
)

// GetPendingRequestsAPI lists incoming follow requests waiting for the
// logged-in user's approval (i.e. their profile is private).
func GetPendingRequestsAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID, err := auth.GetUserIDFromSession(r)
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, "You must be logged in")
		return
	}

	rows, err := database.Database.Query(`
		SELECT u.id, COALESCE(u.first_name, ''), COALESCE(u.email, ''), COALESCE(u.avatar, '')
		FROM followers f
		JOIN users u ON u.id = f.follower_id
		WHERE f.followee_id = ? AND f.status = 'pending'
	`, userID)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Database error")
		return
	}
	defer rows.Close()

	result := []Follow{}
	for rows.Next() {
		var f Follow
		if err := rows.Scan(&f.FollowerID, &f.FirstName, &f.Email, &f.Avatar); err != nil {
			continue
		}
		result = append(result, f)
	}

	utils.RespondJSON(w, http.StatusOK, result)
}