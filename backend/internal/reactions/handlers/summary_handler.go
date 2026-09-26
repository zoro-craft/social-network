package handlers

import (
	"database/sql"

	"zone/internal/database"
	"zone/internal/reactions/models"
)

func getReactionSummary(postID, userID int) (models.Reaction, error) {
	var resp models.Reaction
	resp.PostID = postID

	err := database.Database.QueryRow(`
    SELECT
        COUNT(CASE WHEN is_like = 1 THEN 1 END),
        COUNT(CASE WHEN is_like = 0 THEN 1 END)
    FROM post_reactions
    WHERE post_id = ?`,
		postID,
	).Scan(&resp.LikeCount, &resp.DislikeCount)
	if err != nil {
		return resp, err
	}

	var userIsLike int
	err = database.Database.QueryRow(
		"SELECT is_like FROM post_reactions WHERE post_id = ? AND user_id = ?",
		postID, userID,
	).Scan(&userIsLike)
	if err == nil {
		if userIsLike == 1 {
			resp.UserReaction = "like"
		} else {
			resp.UserReaction = "dislike"
		}
	} else if err != sql.ErrNoRows {
		return resp, err
	}

	return resp, nil
}
