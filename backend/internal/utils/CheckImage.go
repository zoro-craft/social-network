package utils

import (
	"errors"
	"net/http"
	"path/filepath"
	"strings"
)

const MaxUploadSize = 5 * 1024 * 1024

var AllowedImageExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
}
var ErrInvalidImageType = errors.New("only JPEG, PNG, and GIF images are allowed")

func SaveImage(r *http.Request, dir, urlPrefix string) (string, error) {
	file, header, err := r.FormFile("image")
	if err == http.ErrMissingFile {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !AllowedImageExtensions[ext] {
		return "", ErrInvalidImageType
	}

	filename, err := SaveUploadFile(file, header.Filename, dir)
	if err != nil {
		return "", err
	}
	return urlPrefix + filename, nil
}
