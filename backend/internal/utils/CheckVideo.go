package utils

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	// Videos are orders of magnitude larger than the images SaveImage handles,
	// so they get their own ceiling instead of sharing MaxUploadSize.
	MaxVideoUploadSize = 100 * 1024 * 1024

	// Multipart bodies above this size spill to a temp file on disk instead of
	// being buffered in memory, which keeps a 100MB reel from ballooning the heap.
	MultipartMemory = 8 * 1024 * 1024
)

var AllowedVideoExtensions = map[string]bool{
	".mp4":  true,
	".m4v":  true,
	".mov":  true,
	".webm": true,
}

var ErrInvalidVideoType = errors.New("only MP4, MOV and WebM videos are allowed")

// SaveVideo validates a `fieldName` video upload and stores it under uploads/<dir>,
// returning the public URL prefix + generated filename. An empty string with a nil
// error means the form simply carried no file.
func SaveVideo(r *http.Request, fieldName, dir, urlPrefix string) (string, error) {
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

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !AllowedVideoExtensions[ext] {
		return "", ErrInvalidVideoType
	}

	// A container extension is trivially spoofable, so sniff the magic bytes too.
	// MP4/MOV/M4V are ISO-BMFF, where the "ftyp" box starts at offset 4; WebM
	// (and Matroska) begin with the EBML header 0x1A45DFA3.
	head := make([]byte, 12)
	n, err := io.ReadFull(file, head)
	if err != nil && err != io.ErrUnexpectedEOF && err != io.EOF {
		return "", fmt.Errorf("error reading file header")
	}
	if !looksLikeVideo(head[:n], ext) {
		return "", ErrInvalidVideoType
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", fmt.Errorf("error resetting file pointer")
	}

	newFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	targetDir := filepath.Join(UploadsDir(), dir)
	if err := os.MkdirAll(targetDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %v", err)
	}

	dst, err := os.Create(filepath.Join(targetDir, newFileName))
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		// A half-written file would still be served as a broken reel, so drop it.
		os.Remove(dst.Name())
		return "", fmt.Errorf("failed to save video")
	}

	return urlPrefix + newFileName, nil
}

func looksLikeVideo(head []byte, ext string) bool {
	if ext == ".webm" {
		return len(head) >= 4 &&
			head[0] == 0x1A && head[1] == 0x45 && head[2] == 0xDF && head[3] == 0xA3
	}
	return len(head) >= 8 && string(head[4:8]) == "ftyp"
}
