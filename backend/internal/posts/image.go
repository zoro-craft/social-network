package posts

import (
	"errors"
	"net/http"
	"path/filepath"
	"strings"

	"zone/internal/utils"
)

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
	if !utils.AllowedCommentExtensions[ext] {
		return "", utils.ErrInvalidImageType
	}

	filename, err := utils.SaveUploadFile(file, header.Filename, "./uploads/posts")
	if err != nil {
		return "", err
	}
	return "/uploads/posts/" + filename, nil
}
