package posts

import (
	"errors"
	"net/http"
	"path/filepath"
	"strings"

	"zone/internal/utils"
)

var allowedPostExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
}

var errInvalidImageType = errors.New("only JPEG, PNG, and GIF images are allowed")

// savePostImage saves an optional image/GIF attached to the post.
// Returns an empty string if no file was provided.
func savePostImage(r *http.Request) (string, error) {
	file, header, err := r.FormFile("image")
	if err == http.ErrMissingFile {
		return "", errors.New("Image is requird")
	}
	if err != nil {
		return "", err
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !allowedPostExtensions[ext] {
		return "", errInvalidImageType
	}

	filename, err := utils.SaveUploadFile(file, header.Filename, "./uploads/posts")
	if err != nil {
		return "", err
	}
	return "/uploads/posts/" + filename, nil
}
