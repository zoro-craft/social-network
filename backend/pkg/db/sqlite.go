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
	files, _ := filepath.Glob(pattern)

	for _, file := range files {
		content, _ := os.ReadFile(file)
		if _, err := db.Exec(string(content)); err != nil {
			log.Fatal("Failed to migrate database:", err)
		}
		log.Println("Migrated database with file:", file)
	}
}
