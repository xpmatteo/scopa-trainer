package main

import (
	"log"
	"net/http"
	"os"

	"github.com/xpmatteo/scopa-trainer/pkg/adapters/http/handlers"
	"github.com/xpmatteo/scopa-trainer/pkg/adapters/http/views"
	"github.com/xpmatteo/scopa-trainer/pkg/application"
)

func main() {
	// Initialize the application service
	gameService := application.NewGameService()

	// Create template with functions
	templates := views.ParseTemplates("templates/game.html")

	// Set up routes. ALWAYS USE THE METHOD in the route declaration
	http.HandleFunc("GET /", handlers.NewHandleIndex(gameService, templates))
	http.HandleFunc("POST /new-game", handlers.NewHandleNewGame(gameService))
	http.HandleFunc("POST /review-game", handlers.NewHandleReviewGame())
	http.HandleFunc("POST /select-card", handlers.NewHandleSelectCard(gameService))
	http.HandleFunc("POST /play-card", handlers.NewHandlePlayCard(gameService))
	http.HandleFunc("POST /ai-turn", handlers.NewHandleAITurn(gameService))
	http.HandleFunc("POST /test-game-over", handlers.NewHandleTestGameOver(gameService))

	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Get the port from environment variables or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}

	// Start the server
	log.Printf("Starting server on :%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
