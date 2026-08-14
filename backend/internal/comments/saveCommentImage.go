package comments

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"zone/internal/utils"

	"github.com/gofrs/uuid"
)

func saveCommentImage(r *http.Request) (string, error) {
	file, header, err := r.FormFile("image")
	if err != nil {
		return "", nil
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !utils.AllowedCommentExtensions[ext] {
		return "", utils.ErrInvalidImageType
	}

	uploadDir := "./uploads/comments"
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

	return "/uploads/comments/" + filename, nil
}