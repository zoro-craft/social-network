package utils

import (
	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
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

func SaveImageField(r *http.Request, fieldName string, dir string, urlPrefix string) (string, error) {

	file, header, err := r.FormFile(fieldName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))

	if !AllowedImageExtensions[ext] {
		return "", ErrInvalidImageType
	}

	// first check - Magic bytes
	buff := make([]byte, 512)
	if _, err := file.Read(buff); err != nil && err != io.EOF {
		return "", fmt.Errorf("error reading file header")
	}
	fileType := http.DetectContentType(buff)
	if fileType != "image/jpeg" && fileType != "image/png" && fileType != "image/gif" {
		return "", ErrInvalidImageType
	}

	// start index reading into 0
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", fmt.Errorf("error resetting file pointer")
	}

	// second check -- size of image
	config, _, err := image.DecodeConfig(file)
	if err != nil {
		return "", fmt.Errorf("invalid image metadata")
	}
	if config.Width > 4000 || config.Height > 4000 {
		return "", fmt.Errorf("image dimensions too large (max 4000x4000 pixels)")
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", fmt.Errorf("error resetting file pointer")
	}

	// last check -- rendering image
	img, _, err := image.Decode(file)
	if err != nil {
		return "", fmt.Errorf("invalid or corrupted image file")
	}

	newFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	fullPath := filepath.Join(".", "uploads", dir, newFileName)

	dst, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// encoding image into origin ext
	switch ext {
	case ".jpg", ".jpeg":
		err = jpeg.Encode(dst, img, nil)
	case ".png":
		err = png.Encode(dst, img)
	case ".gif":
		err = gif.Encode(dst, img, nil)
	}

	if err != nil {
		return "", fmt.Errorf("failed to process and save image")
	}

	return urlPrefix + newFileName, nil
}

func ServeUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	rel := strings.TrimPrefix(r.URL.Path, "/uploads")
	http.ServeFile(w, r, filepath.Join(UploadsDir(), rel))
}
