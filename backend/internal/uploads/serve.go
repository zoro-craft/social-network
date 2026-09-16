package uploads

import (
	"database/sql"
	"net/http"
	"path/filepath"
	"strings"

	"zone/internal/auth"
	"zone/internal/database"
	"zone/internal/posts"
	"zone/internal/utils"
)

func ServeImage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := auth.GetUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var postID int
	var dir string

	switch {
	case strings.HasPrefix(r.URL.Path, "/uploads/posts/"):
		dir, postID, err = postImagePostID("/uploads/posts/", r.URL.Path)
	case strings.HasPrefix(r.URL.Path, "/uploads/comments/"):
		dir, postID, err = commentImagePostID("/uploads/comments/", r.URL.Path)
	default:
		http.NotFound(w, r)
		return
	}

	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if !posts.CanViewPost(postID, userID) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	http.ServeFile(w, r, filepath.Join(dir, filepath.Base(r.URL.Path)))
}

func postImagePostID(imagePrefix, urlPath string) (string, int, error) {
	var postID int
	err := database.Database.QueryRow(
		"SELECT id FROM posts WHERE image = ?",
		imagePrefix+filepath.Base(urlPath),
	).Scan(&postID)
	return filepath.Join(utils.UploadsDir(), "posts"), postID, err
}

func commentImagePostID(imagePrefix, urlPath string) (string, int, error) {
	var postID int
	err := database.Database.QueryRow(
		"SELECT post_id FROM comments WHERE image = ?",
		imagePrefix+filepath.Base(urlPath),
	).Scan(&postID)
	return filepath.Join(utils.UploadsDir(), "comments"), postID, err
}
