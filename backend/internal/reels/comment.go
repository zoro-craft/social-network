package reels

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"zone/internal/auth"
	"zone/internal/database"
	"zone/internal/utils"
)

const maxCommentLength = 1000

const commentSelect = `
	SELECT c.id, c.reel_id, c.user_id, c.content,
	       COALESCE(u.first_name, ''), COALESCE(u.avatar, ''), c.created_at
	FROM reel_comments c
	JOIN users u ON c.user_id = u.id
`

func GetReelCommentsAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	if _, err := auth.GetUserID(r); err != nil {
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

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 50 {
		limit = 20
	}

	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if offset < 0 {
		offset = 0
	}

	query := commentSelect + " WHERE c.reel_id = ?"
	args := []any{reelID}

	beforeStr := r.URL.Query().Get("before")
	if beforeStr != "" {
		beforeID, err := strconv.Atoi(beforeStr)
		if err != nil || beforeID <= 0 {
			utils.HandleError(w, http.StatusBadRequest, "Invalid before id")
			return
		}
		query += " AND c.id < ?"
		args = append(args, beforeID)
	}

	query += " ORDER BY c.id DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := database.Database.Query(query, args...)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Database error loading comments")
		return
	}
	defer rows.Close()

	result := []ReelComment{}
	for rows.Next() {
		var c ReelComment
		if err := rows.Scan(
			&c.ID, &c.ReelID, &c.UserID, &c.Content, &c.FirstName, &c.Avatar, &c.CreatedAt,
		); err != nil {
			utils.HandleError(w, http.StatusInternalServerError, "Error reading comments")
			return
		}
		result = append(result, c)
	}

	if err := rows.Err(); err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Error reading comments")
		return
	}

	if beforeStr == "" {
		var maxID int
		database.Database.QueryRow(
			"SELECT COALESCE(MAX(id), 0) FROM reel_comments WHERE reel_id = ?", reelID,
		).Scan(&maxID)
		w.Header().Set("X-Max-Comment-Id", strconv.Itoa(maxID))
	}

	utils.RespondJSON(w, http.StatusOK, result)
}

func CreateReelCommentAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID, err := auth.GetUserID(r)
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, "You must be logged in to comment")
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

	// Comments are text-only, so a small body cap is enough here. ParseMultipartForm
	// also parses urlencoded bodies, so clients may send either encoding.
	r.Body = http.MaxBytesReader(w, r.Body, 32*1024)
	if err := r.ParseMultipartForm(32 * 1024); err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Invalid form data")
		return
	}

	content := strings.TrimSpace(r.FormValue("content"))
	if content == "" {
		utils.HandleError(w, http.StatusBadRequest, "Content is required")
		return
	}
	if len(content) > maxCommentLength {
		utils.HandleError(w, http.StatusBadRequest, "Content too long")
		return
	}

	result, err := database.Database.Exec(
		"INSERT INTO reel_comments (reel_id, user_id, content, created_at) VALUES (?, ?, ?, ?)",
		reelID, userID, content, time.Now().UnixMilli(),
	)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to create comment")
		return
	}

	commentID, err := result.LastInsertId()
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Could not get new comment ID")
		return
	}

	var c ReelComment
	err = database.Database.QueryRow(commentSelect+" WHERE c.id = ?", commentID).Scan(
		&c.ID, &c.ReelID, &c.UserID, &c.Content, &c.FirstName, &c.Avatar, &c.CreatedAt,
	)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to load created comment")
		return
	}

	utils.RespondJSON(w, http.StatusCreated, c)
}
