package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"zone/internal/database"
	"zone/internal/routing"
	"zone/internal/seeddata"
	"zone/internal/websocket"
)

func main() {
	// Initialize database
	if err := database.Init(); err != nil {
		log.Fatal(err)
	}
	defer database.Database.Close()

	// Seed demo data on first run. Existing accounts are never overwritten.
	if seedDemoEnabled() {
		if _, err := seeddata.Run(database.Database, false); err != nil {
			log.Printf("Seed failed: %v", err)
		}
	} else {
		log.Println("Seed: auto-seed disabled (SEED_DEMO=false)")
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

// seedDemoEnabled defaults to true so a fresh database is usable immediately.
// Set SEED_DEMO=false when a deployment must not create demo accounts.
func seedDemoEnabled() bool {
	value, ok := os.LookupEnv("SEED_DEMO")
	if !ok {
		return true
	}

	switch strings.ToLower(strings.TrimSpace(value)) {
	case "false", "0", "no", "off":
		return false
	default:
		return true
	}
}
