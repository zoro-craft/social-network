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

const (
	MaxUploadSize = 5 * 1024 * 1024

	// Ceiling for a whole multipart body, not just the file in it. The difference
	// covers the envelope (boundaries, part headers, text fields), so an image of
	// exactly MaxUploadSize still fits instead of being rejected as "too large".
	MaxImageFormSize = MaxUploadSize + (1 << 20)
)

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
		if err == http.ErrMissingFile {
			return "", nil
		}
		return "", err
	}
	defer file.Close()
	if header.Size == 0 || header.Filename == "" {
		return "", nil
	}
	if header.Size > MaxUploadSize {
		return "", fmt.Errorf("image is too large (max 5 MB)")
	}
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
	// GIFs need special handling: image.Decode only reads the first frame,
	// which silently flattens animated GIFs into a static image.
	var gifImg *gif.GIF
	var img image.Image
	if ext == ".gif" {
		gifImg, err = gif.DecodeAll(file)
		if err != nil {
			return "", fmt.Errorf("invalid or corrupted image file")
		}
	} else {
		img, _, err = image.Decode(file)
		if err != nil {
			return "", fmt.Errorf("invalid or corrupted image file")
		}
	}

	newFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)

	baseUploadDir := UploadsDir()
	targetDir := filepath.Join(baseUploadDir, dir)
	fullPath := filepath.Join(targetDir, newFileName)

	if err := os.MkdirAll(targetDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %v", err)
	}

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
		err = gif.EncodeAll(dst, gifImg)
	}

	if err != nil {
		return "", fmt.Errorf("failed to process and save image")
	}

	return urlPrefix + newFileName, nil
}

// ServeUpload serves a single file from one subdirectory of the uploads root,
// e.g. /uploads/reels/clip.mp4. Only "<subdir>/<name>" is accepted: r.URL.Path is
// already percent-decoded at this point, so /uploads/reels/%2e%2e/%2e%2e/data/db
// arrives as "reels/../../data/db" and ServeMux does not reject it, since it only
// canonicalises the still-escaped form of the path.
func ServeUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	subDir, name, found := strings.Cut(strings.TrimPrefix(r.URL.Path, "/uploads/"), "/")
	if !found || subDir == "" || name == "" || strings.Contains(name, "/") {
		http.NotFound(w, r)
		return
	}

	http.ServeFile(w, r, filepath.Join(UploadsDir(), subDir, name))
}
