package main

import (
	"log"
	"net/http"

	"social-network/internal/auth"
	"social-network/pkg/db"
)

func main() {
	database := db.InitDB("./social_network.db")
	defer database.Close()

	db.MigrateDB(database, "./pkg/db/migrations/sqlite/*.up.sql")

	http.HandleFunc("/api/register", auth.RegisterHandler)
	http.HandleFunc("/api/login", auth.LoginHandler)

	log.Println("Server is running on http://localhost:8080")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
