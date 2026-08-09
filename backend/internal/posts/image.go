package posts

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofrs/uuid"
)

// savePostImage saves an optional image/GIF attached to the post.
// Returns an empty string if no file was provided.
func savePostImage(r *http.Request) (string, error) {
	file, header, err := r.FormFile("image")
	if err != nil {
		return "", errors.New("an image is required")
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !allowedPostExtensions[ext] {
		return "", errInvalidImageType
	}

	uploadDir := "./uploads/posts"
	if err := os.MkdirAll(uploadDir, 0o755); err != nil {
		return "", err
	}

	newUUID, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	filename := newUUID.String() + ext
	fullPath := filepath.Join(uploadDir, filename)
	dst, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", err
	}

	return "/uploads/posts/" + filename, nil
}