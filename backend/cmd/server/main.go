package main

import (
	"log"
	"net/http"

	"zone/internal/database"
	"zone/internal/routing"
)

func main() {
	if err := database.Init(); err != nil {
		log.Fatal(err)
	}

	log.Println("Server running on :8080")
	err := http.ListenAndServe(":8080", routing.NewRouter())
	if (err != nil) {
		log.Fatal(err)
	}
}