package uploads

import (
	"path/filepath"

	"zone/internal/database"
)

func postImagePostID(imagePrefix, urlPath string) (string, int, error) {
	var postID int
	err := database.Database.QueryRow(
		"SELECT id FROM posts WHERE image = ?",
		imagePrefix+filepath.Base(urlPath),
	).Scan(&postID)
	return "./uploads/posts", postID, err
}

func commentImagePostID(imagePrefix, urlPath string) (string, int, error) {
	var postID int
	err := database.Database.QueryRow(
		"SELECT post_id FROM comments WHERE image = ?",
		imagePrefix+filepath.Base(urlPath),
	).Scan(&postID)
	return "./uploads/comments", postID, err
}
