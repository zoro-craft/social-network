package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofrs/uuid"
)

var allowedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
}

const maxUploadSize = 20 << 20 // 20MB

func UploadAvatar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		HandleError(w, http.StatusBadRequest, "File too large (max 20MB)")
		return
	}

	file, header, err := r.FormFile("avatar")
	if err != nil {
		HandleError(w, http.StatusBadRequest, "No avatar file provided")
		return
	}
	defer file.Close()

	if !allowedExtensions[fileExt(header.Filename)] {
		HandleError(w, http.StatusBadRequest, "Only JPEG, PNG, and GIF files are allowed")
		return
	}

	filename, err := SaveUploadFile(file, header.Filename, "./uploads/avatars")
	if err != nil {
		HandleError(w, http.StatusInternalServerError, "Could not save file")
		return
	}

	RespondJSON(w, http.StatusOK, map[string]string{
		"path": fmt.Sprintf("/uploads/avatars/%s", filename),
	})
}

func SaveUploadFile(src io.Reader, originalName, dir string) (string, error) {
	//creat folder
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}

	newUUID, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	filename := newUUID.String() + fileExt(originalName)

	dst, err := os.Create(filepath.Join(dir, filename))
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", err
	}
	return filename, nil
}

func fileExt(name string) string {
	return strings.ToLower(filepath.Ext(name))
}
