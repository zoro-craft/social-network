package posts

import (
	"net/http"
	"strconv"

	"zone/internal/auth"
	"zone/internal/database"
	"zone/internal/utils"
)

// visibleCondition returns a SQL condition true when a post is visible to
// viewerID. Needs viewerID passed 3 times as query args.
const visibleCondition = `(
	p.privacy = 'public'
	OR p.user_id = ?
	OR (p.privacy = 'almost_private' AND EXISTS (
		SELECT 1 FROM followers f
		WHERE f.follower_id = ? AND f.followee_id = p.user_id AND f.status = 'accepted'
	))
	OR (p.privacy = 'private' AND EXISTS (
		SELECT 1 FROM post_viewers pv WHERE pv.post_id = p.id AND pv.user_id = ?
	))
)`

func GetPostsAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	viewerID, _ := auth.GetUserIDFromSession(r)

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if offset < 0 {
		offset = 0
	}

	query := `
		SELECT p.id, p.title, p.content, COALESCE(p.image, ''), p.privacy,
		       p.user_id, COALESCE(u.nickname, ''), p.created_at
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE ` + visibleCondition + `
		ORDER BY p.id DESC LIMIT ? OFFSET ?
	`

	rows, err := database.Database.Query(query, viewerID, viewerID, viewerID, limit, offset)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Database error loading posts")
		return
	}
	defer rows.Close()

	posts := []Post{}
	for rows.Next() {
		var p Post
		if rows.Scan(&p.ID, &p.Title, &p.Content, &p.Image, &p.Privacy, &p.UserID, &p.Nickname, &p.CreatedAt) == nil {
			posts = append(posts, p)
		}
	}

	utils.RespondJSON(w, http.StatusOK, posts)
}

// GetPostsByUser returns userID's posts visible to viewerID.
func GetPostsByUser(userID, viewerID int) ([]Post, error) {
	query := `
		SELECT p.id, p.title, p.content, COALESCE(p.image, ''), p.privacy,
		       p.user_id, COALESCE(u.nickname, ''), p.created_at
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.user_id = ? AND ` + visibleCondition + `
		ORDER BY p.id DESC
	`

	rows, err := database.Database.Query(query, userID, viewerID, viewerID, viewerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []Post{}
	for rows.Next() {
		var p Post
		if rows.Scan(&p.ID, &p.Title, &p.Content, &p.Image, &p.Privacy, &p.UserID, &p.Nickname, &p.CreatedAt) == nil {
			result = append(result, p)
		}
	}
	return result, rows.Err()
}