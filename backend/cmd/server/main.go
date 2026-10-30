package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"zone/internal/database"
	"zone/internal/routing"
	"zone/internal/websocket"
)

func main() {
	// Initialize database
	if err := database.Init(); err != nil {
		log.Fatal(err)
	}

	// Fill the database with demo accounts and content on first startup.
	if err := database.Seed(); err != nil {
		log.Fatal(err)
	}

	// Start the realtime hub goroutine
	websocket.WirePresence(websocket.DefaultHub)
	go websocket.DefaultHub.Run()

	// Create central router
	router := routing.NewRouter()
	server := &http.Server{Addr: ":8080", Handler: router}

	go func() {
		log.Println("Server running on http://localhost:8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// Graceful shutdown: on SIGINT/SIGTERM close the hub and stop the server.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	websocket.DefaultHub.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Println("Server shutdown:", err)
	}
	log.Println("Server stopped")
}
