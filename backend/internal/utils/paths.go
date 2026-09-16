package utils

import "os"

func UploadsDir() string {
	if p := os.Getenv("UPLOADS_PATH"); p != "" {
		return p
	}
	return "./uploads"
}