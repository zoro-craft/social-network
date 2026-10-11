package reels

import (
	"database/sql"
	"net/http"
	"strconv"

	"zone/internal/auth"
	"zone/internal/database"
	"zone/internal/utils"
)

// reelSelect carries the per-viewer flags (liked/saved) as its two leading args,
// so every caller must pass the viewer id twice before its own arguments.
const reelSelect = `
	SELECT
		r.id,
		r.user_id,
		u.full_name,
		COALESCE(u.avatar, ''),
		r.video,
		r.caption,
		r.created_at,
		(SELECT COUNT(*) FROM reel_likes    rl WHERE rl.reel_id = r.id)  AS like_count,
		(SELECT COUNT(*) FROM reel_comments rc WHERE rc.reel_id = r.id)  AS comment_count,
		(SELECT COUNT(*) FROM reel_views    rv WHERE rv.reel_id = r.id)  AS view_count,
		(SELECT COUNT(*) FROM reel_saves    rs WHERE rs.reel_id = r.id)  AS save_count,
		(SELECT COUNT(*) FROM reel_shares   rh WHERE rh.reel_id = r.id)  AS share_count,
		CASE WHEN EXISTS(SELECT 1 FROM reel_likes rl WHERE rl.reel_id = r.id AND rl.user_id = ?) THEN 1 ELSE 0 END,
		CASE WHEN EXISTS(SELECT 1 FROM reel_saves rs WHERE rs.reel_id = r.id AND rs.user_id = ?) THEN 1 ELSE 0 END
	FROM reels r
	JOIN users u ON r.user_id = u.id
`

// scanner is satisfied by both *sql.Row and *sql.Rows, so one scan helper covers
// the list and single-reel paths.
type scanner interface {
	Scan(dest ...any) error
}

func scanReel(row scanner) (Reel, error) {
	var reel Reel
	var liked, saved bool

	err := row.Scan(
		&reel.ID, &reel.UserID, &reel.FullName, &reel.Avatar,
		&reel.Video, &reel.Caption, &reel.CreatedAt,
		&reel.LikeCount, &reel.CommentCount, &reel.ViewCount,
		&reel.SaveCount, &reel.ShareCount,
		&liked, &saved,
	)
	if err != nil {
		return Reel{}, err
	}

	reel.Liked = liked
	reel.Saved = saved
	return reel, nil
}

// reelIDFromPath reads the positive {id} path value, returning 0 when absent or malformed.
func reelIDFromPath(r *http.Request) int {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		return 0
	}
	return id
}

// Reels are addressable by any signed-in user, so visibility is just existence.
func reelExists(reelID int) bool {
	var exists bool
	database.Database.QueryRow("SELECT EXISTS(SELECT 1 FROM reels WHERE id = ?)", reelID).Scan(&exists)
	return exists
}

// GetReel loads one reel with the viewer's liked/saved flags.
func GetReel(reelID, viewerID int) (Reel, error) {
	return scanReel(database.Database.QueryRow(reelSelect+" WHERE r.id = ?", viewerID, viewerID, reelID))
}

// pageParams reads the shared cursor pagination query string (limit/offset/before).
func pageParams(r *http.Request) (limit, offset, beforeID int, hasBefore bool, err error) {
	limit, _ = strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 50 {
		limit = 10
	}

	offset, _ = strconv.Atoi(r.URL.Query().Get("offset"))
	if offset < 0 {
		offset = 0
	}

	beforeStr := r.URL.Query().Get("before")
	if beforeStr == "" {
		return limit, offset, 0, false, nil
	}

	beforeID, convErr := strconv.Atoi(beforeStr)
	if convErr != nil || beforeID <= 0 {
		return 0, 0, 0, false, strconv.ErrSyntax
	}
	return limit, offset, beforeID, true, nil
}

// listReels runs reelSelect with a caller-supplied WHERE/ORDER and returns the rows.
func listReels(query string, args ...any) ([]Reel, error) {
	rows, err := database.Database.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []Reel{}
	for rows.Next() {
		reel, err := scanReel(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, reel)
	}

	return result, rows.Err()
}

func GetReelsAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	viewerID, err := auth.GetUserID(r)
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, "You must be logged in")
		return
	}

	limit, offset, beforeID, hasBefore, err := pageParams(r)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Invalid before id")
		return
	}

	query := reelSelect
	args := []any{viewerID, viewerID}

	if hasBefore {
		query += " WHERE r.id < ?"
		args = append(args, beforeID)
	}

	query += " ORDER BY r.id DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	result, err := listReels(query, args...)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Database error loading reels")
		return
	}

	// The first page pins a snapshot ceiling so later pages stay consistent while
	// new reels keep arriving.
	if !hasBefore {
		var maxID int
		if err := database.Database.QueryRow("SELECT COALESCE(MAX(id), 0) FROM reels").Scan(&maxID); err == nil {
			w.Header().Set("X-Max-Reel-Id", strconv.Itoa(maxID))
		}
	}

	utils.RespondJSON(w, http.StatusOK, result)
}

func GetReelAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	viewerID, err := auth.GetUserID(r)
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, "You must be logged in")
		return
	}

	reelID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || reelID <= 0 {
		utils.HandleError(w, http.StatusBadRequest, "Invalid reel id")
		return
	}

	reel, err := GetReel(reelID, viewerID)
	if err == sql.ErrNoRows {
		utils.HandleError(w, http.StatusNotFound, "Reel not found")
		return
	}
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to load reel")
		return
	}

	utils.RespondJSON(w, http.StatusOK, reel)
}

// GetSavedReelsAPI backs the "Saved" tab: the viewer's bookmarked reels, newest first.
func GetSavedReelsAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	viewerID, err := auth.GetUserID(r)
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, "You must be logged in")
		return
	}

	limit, offset, _, _, err := pageParams(r)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Invalid before id")
		return
	}

	query := reelSelect + " WHERE EXISTS(SELECT 1 FROM reel_saves s WHERE s.reel_id = r.id AND s.user_id = ?)" +
		" ORDER BY r.id DESC LIMIT ? OFFSET ?"

	result, err := listReels(query, viewerID, viewerID, viewerID, limit, offset)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Database error loading saved reels")
		return
	}

	utils.RespondJSON(w, http.StatusOK, result)
}

// getReelStats re-reads every counter plus the viewer's own flags after a mutation.
func getReelStats(reelID, viewerID int) (ReelStats, error) {
	stats := ReelStats{ReelID: reelID}
	var liked, saved bool

	err := database.Database.QueryRow(`
		SELECT
			(SELECT COUNT(*) FROM reel_likes    WHERE reel_id = ?),
			(SELECT COUNT(*) FROM reel_comments WHERE reel_id = ?),
			(SELECT COUNT(*) FROM reel_views    WHERE reel_id = ?),
			(SELECT COUNT(*) FROM reel_saves    WHERE reel_id = ?),
			(SELECT COUNT(*) FROM reel_shares   WHERE reel_id = ?),
			EXISTS(SELECT 1 FROM reel_likes WHERE reel_id = ? AND user_id = ?),
			EXISTS(SELECT 1 FROM reel_saves WHERE reel_id = ? AND user_id = ?)
	`,
		reelID, reelID, reelID, reelID, reelID,
		reelID, viewerID,
		reelID, viewerID,
	).Scan(
		&stats.LikeCount, &stats.CommentCount, &stats.ViewCount,
		&stats.SaveCount, &stats.ShareCount,
		&liked, &saved,
	)
	if err != nil {
		return stats, err
	}

	stats.Liked = liked
	stats.Saved = saved
	return stats, nil
}

// respondWithStats is the shared tail of like/save/share/view: refresh counters and reply.
func respondWithStats(w http.ResponseWriter, reelID, viewerID int) {
	stats, err := getReelStats(reelID, viewerID)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Could not load reel counters")
		return
	}

	utils.RespondJSON(w, http.StatusOK, stats)
}
