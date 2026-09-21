package database

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "github.com/mattn/go-sqlite3"
)

var Database *sql.DB

// Init sets up the SQLite database connection and applies all pending migrations to ensure the database schema is up to date.
func Init() error {
	const baseDirName = "./internal/database"
	const dbFileName = "social-network.db"

	dbFile := baseDirName + "/" + dbFileName
	if p := os.Getenv("DB_PATH"); p != "" {
		dbFile = p
	}

	var err error
	Database, err = sql.Open("sqlite3", dbFile+"?_foreign_keys=on")
	if err != nil {
		return fmt.Errorf("can't open/create social-network.db: %v", err)
	}

	if err := Database.Ping(); err != nil {
		return fmt.Errorf("can't connect to database: %v", err)
	}
	//  create a new database driver instance for the given database connection. whit the default configuration
	driver, err := sqlite3.WithInstance(Database, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("can't create migration driver: %v", err)
	}

	migrationsPath := "file://" + baseDirName + "/migrations/sqlite"
	// creat the manger ho is connected to the database and the migration files, and run the migrations
	m, err := migrate.NewWithDatabaseInstance(migrationsPath, "sqlite3", driver)
	if err != nil {
		return fmt.Errorf("can't create migrate instance: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration failed: %v", err)
	}

	return nil
}