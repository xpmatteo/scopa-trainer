package main

import (
	"log"
	"net/http"

	"github.com/xpmatteo/scopa-trainer/pkg/adapters/http/handlers"
	"github.com/xpmatteo/scopa-trainer/pkg/application"
)

func main() {
	// Initialize the application service
	gameService := application.NewGameService()

	// Initialize the HTTP handler
	handler, err := handlers.NewHandler(gameService)
	if err != nil {
		log.Fatalf("Failed to initialize handler: %v", err)
	}

	// Set up routes
	http.HandleFunc("GET /", handler.HandleIndex)
	http.HandleFunc("POST /new-game", handler.HandleNewGame)
	// Start the server
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
