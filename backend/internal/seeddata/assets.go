package seeddata

import (
	"embed"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"zone/internal/utils"
)

//go:embed assets
var demoAssets embed.FS

// assetURL turns a path inside assets/ into the public URL the handlers serve.
// copyAssets must have written the matching file under uploads/ first.
func assetURL(rel string) string {
	if rel == "" {
		return ""
	}
	return "/uploads/" + rel
}

// nullableAssetURL keeps optional media columns NULL when no asset is assigned.
func nullableAssetURL(rel string) any {
	if rel == "" {
		return nil
	}
	return assetURL(rel)
}

// copyAssets writes the embedded demo media into the uploads tree the server
// serves. The stable seed-* names mean a reseed overwrites the same files
// instead of accumulating timestamped copies.
func copyAssets() error {
	return fs.WalkDir(demoAssets, "assets", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		data, err := demoAssets.ReadFile(path)
		if err != nil {
			return err
		}

		rel := strings.TrimPrefix(path, "assets/")
		dst := filepath.Join(utils.UploadsDir(), filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			return err
		}
		return os.WriteFile(dst, data, 0o644)
	})
}
