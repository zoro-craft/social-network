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
	return SaveImageField(r, "image", dir, urlPrefix)
}

// SaveImageField saves an uploaded image read from the given multipart field
// name. Posts and comments send their file as "image"; avatar uploads
// (register + edit profile) send theirs as "avatar".
func SaveImageField(r *http.Request, field, dir, urlPrefix string) (string, error) {
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

	filename, err := SaveUploadFile(file, header.Filename, dir)
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

	http.ServeFile(w, r, filepath.Join(".", r.URL.Path))
}
