package db

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
)

func InitDB(dbPath string) *sql.DB {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}
	log.Println("Connected to database")
	return db
}

func MigrateDB(db *sql.DB, pattern string) {
	files, _ := filepath.Glob(pattern)

	for _, file := range files {
		content, _ := os.ReadFile(file)
		if _, err := db.Exec(string(content)); err != nil {
			log.Fatal("Failed to migrate database:", err)
		}
		log.Println("Migrated database with file:", file)
	}
}
