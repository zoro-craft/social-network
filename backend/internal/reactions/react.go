package reactions

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"zone/internal/auth"
	"zone/internal/database"
	"zone/internal/posts"
	"zone/internal/utils"
)

func ReactToPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID, err := auth.GetUserID(r)
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req Reaction
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Type != "like" && req.Type != "dislike" {
		utils.HandleError(w, http.StatusBadRequest, "type must be 'like' or 'dislike'")
		return
	}

	if !posts.CanViewPost(req.PostID, userID) {
		utils.HandleError(w, http.StatusBadRequest, "Invalid post")
		return
	}

	newIsLike := 0
	if req.Type == "like" {
		newIsLike = 1
	}

	var existingIsLike int
	err = database.Database.QueryRow(
		"SELECT is_like FROM post_reactions WHERE user_id = ? AND post_id = ?",
		userID, req.PostID,
	).Scan(&existingIsLike)

	switch {
	case err == sql.ErrNoRows:
		_, err = database.Database.Exec(
			"INSERT INTO post_reactions (user_id, post_id, is_like) VALUES (?, ?, ?)",
			userID, req.PostID, newIsLike,
		)
	case err == nil && existingIsLike == newIsLike:
		_, err = database.Database.Exec(
			"DELETE FROM post_reactions WHERE user_id = ? AND post_id = ?",
			userID, req.PostID,
		)
	case err == nil:
		_, err = database.Database.Exec(
			"UPDATE post_reactions SET is_like = ? WHERE user_id = ? AND post_id = ?",
			newIsLike, userID, req.PostID,
		)
	}
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Could not save reaction")
		return
	}

	resp, err := getReactionSummary(req.PostID, userID)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Could not load reaction counts")
		return
	}

	utils.RespondJSON(w, http.StatusOK, resp)
}
