package utils

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofrs/uuid"
)

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
