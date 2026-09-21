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

func SaveImage(r *http.Request, subDir, urlPrefix string) (string, error) {
	return SaveImageField(r, "image", subDir, urlPrefix)
}

func SaveImageField(r *http.Request, field, subDir, urlPrefix string) (string, error) {
	file, header, err := r.FormFile(field)
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

	filename, err := SaveUploadFile(file, header.Filename, filepath.Join(UploadsDir(), subDir))
	if err != nil {
		return "", err
	}
	return urlPrefix + filename, nil
}

func ServeUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	rel := strings.TrimPrefix(r.URL.Path, "/uploads")
	http.ServeFile(w, r, filepath.Join(UploadsDir(), rel))
}