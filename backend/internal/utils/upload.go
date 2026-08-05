package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gofrs/uuid"
)

var allowedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
}

const maxUploadSize = 5 << 20 // 5MB

func UploadAvatar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		HandleError(w, http.StatusBadRequest, "File too large (max 5MB)")
		return
	}

	file, header, err := r.FormFile("avatar")
	if err != nil {
		HandleError(w, http.StatusBadRequest, "No avatar file provided")
		return
	}
	defer file.Close()

	ext := filepath.Ext(header.Filename)
	if !allowedExtensions[ext] {
		HandleError(w, http.StatusBadRequest, "Only JPEG, PNG, and GIF files are allowed")
		return
	}

	uploadDir := "./uploads/avatars"
	if err := os.MkdirAll(uploadDir, 0o755); err != nil {
		HandleError(w, http.StatusInternalServerError, "Could not create upload directory")
		return
	}

	newUUID, err := uuid.NewV4()
	if err != nil {
		HandleError(w, http.StatusInternalServerError, "Server error")
		return
	}
	filename := newUUID.String() + ext
	fullPath := filepath.Join(uploadDir, filename)

	dst, err := os.Create(fullPath)
	if err != nil {
		HandleError(w, http.StatusInternalServerError, "Could not save file")
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		HandleError(w, http.StatusInternalServerError, "Could not save file")
		return
	}

	RespondJSON(w, http.StatusOK, map[string]string{
		"path": fmt.Sprintf("/uploads/avatars/%s", filename),
	})
}