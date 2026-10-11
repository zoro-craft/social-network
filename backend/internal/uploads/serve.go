package uploads

import (
	"database/sql"
	"net/http"
	"path/filepath"
	"strings"

	"zone/internal/auth"
	"zone/internal/database"
	"zone/internal/posts/repositories"
	"zone/internal/stories"
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

	// Story images have their own audience rules, so they bypass the post switch.
	if strings.HasPrefix(r.URL.Path, "/uploads/stories/") {
		serveStoryImage(w, r, userID)
		return
	}

	if strings.HasPrefix(r.URL.Path, "/uploads/groups/") {
		if !groupImageExists(r.URL.Path) {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, filepath.Join(utils.UploadsDir(), "groups", filepath.Base(r.URL.Path)))
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

	if !repositories.CanViewPost(postID, userID) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	http.ServeFile(w, r, filepath.Join(dir, filepath.Base(r.URL.Path)))
}

// serveStoryImage serves a story only to its author and their accepted followers,
// and only while the story is still unexpired.
func serveStoryImage(w http.ResponseWriter, r *http.Request, userID int) {
	var storyID int
	err := database.Database.QueryRow(
		"SELECT id FROM stories WHERE image = ?",
		"/uploads/stories/"+filepath.Base(r.URL.Path),
	).Scan(&storyID)

	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if !stories.CanViewStory(storyID, userID) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	http.ServeFile(w, r, filepath.Join(utils.UploadsDir(), "stories", filepath.Base(r.URL.Path)))
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

func groupImageExists(urlPath string) bool {
	var groupID int
	err := database.Database.QueryRow(
		"SELECT id FROM user_groups WHERE image = ?",
		"/uploads/groups/"+filepath.Base(urlPath),
	).Scan(&groupID)
	return err == nil
}
