// Command seed fills the SQLite database with the shared demo dataset.
//
// It reuses internal/database so it always targets the same schema as the
// server: database.Init applies any pending migrations, and DB_PATH selects
// the file, so the same binary works locally and inside the container.
//
// Usage:
//
//	go run ./cmd/seed              # wipe app data, then seed
//	go run ./cmd/seed -db /tmp/x.db
//	go run ./cmd/seed -wipe=false  # seed only if the database is empty
package main

import (
	"flag"
	"log"
	"os"

	"zone/internal/database"
	"zone/internal/seeddata"
)

func main() {
	log.SetFlags(0)

	dbPath := flag.String("db", "", "database file to seed (defaults to DB_PATH, then the local dev database)")
	wipe := flag.Bool("wipe", true, "delete existing application data before seeding")
	flag.Parse()

	if *dbPath != "" {
		if err := os.Setenv("DB_PATH", *dbPath); err != nil {
			log.Fatalf("set DB_PATH: %v", err)
		}
	}

	if err := database.Init(); err != nil {
		log.Fatalf("init database: %v", err)
	}
	defer database.Database.Close()

	result, err := seeddata.Run(database.Database, *wipe)
	if err != nil {
		log.Fatalf("seed: %v", err)
	}
	if !result.Seeded {
		return
	}

	printReport(result)
}
