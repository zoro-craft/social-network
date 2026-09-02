package main

import (
	"log"
	"net/http"

	"zone/internal/database"
	"zone/internal/routing"
	"zone/pkg/db"
)

func main() {
	// Initialize Posts/Media database
	if err := database.Init(); err != nil {
		log.Fatal(err)
	}

	// Initialize Auth database
	authDB := db.InitDB("./social_network.db")
	defer authDB.Close()
	db.MigrateDB(authDB, "./pkg/db/migrations/sqlite/*.up.sql")

	// Create central router
	router := routing.NewRouter()

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
