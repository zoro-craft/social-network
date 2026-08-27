package main

import (
	"log"

	"social-network/pkg/db"
)

func main() {
	log.Println("🚀 Starting Social Network Server...")

	database := db.InitDB("./social_network.db")
	defer database.Close()

	db.MigrateDB(database, "./pkg/db/migrations/sqlite/*.up.sql")

	log.Println("✅ Database Setup Complete! Server is ready.")

}
