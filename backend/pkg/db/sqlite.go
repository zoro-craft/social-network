package db

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB(dbPath string) *sql.DB {
	var err error
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}
	if err = DB.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}
	log.Println("Connected to database")
	return DB
}

func MigrateDB(db *sql.DB, pattern string) {
	files, err := filepath.Glob(pattern)
	if err != nil || len(files) == 0 {
		fallbackPattern := filepath.Join("backend", pattern)
		fallbackFiles, fallbackErr := filepath.Glob(fallbackPattern)
		if fallbackErr == nil && len(fallbackFiles) > 0 {
			files = fallbackFiles
		} else {
			log.Printf("Warning: No migration files found for pattern: %s (or fallback: %s)\n", pattern, fallbackPattern)
			return
		}
	}

	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			log.Printf("Failed to read migration file %s: %v\n", file, err)
			continue
		}
		if _, err := db.Exec(string(content)); err != nil {
			log.Fatalf("Failed to migrate database with file %s: %v\n", file, err)
		}
		log.Println("Migrated database with file:", file)
	}
}
