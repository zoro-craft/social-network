package main

import (
	"log"
	"net/http"

	"zone/internal/database"
	"zone/internal/routing"
)

func main() {
	// Initialize Posts/Media database
	if err := database.Init(); err != nil {
		log.Fatal(err)
	}


	// Create central router
	router := routing.NewRouter()

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
