package main

import (
	"log"
	"net/http"

	"github.com/xpmatteo/scopa-trainer/pkg/adapters/http/handlers"
	"github.com/xpmatteo/scopa-trainer/pkg/adapters/http/views"
	"github.com/xpmatteo/scopa-trainer/pkg/application"
)

func main() {
	// Initialize the application service
	gameService := application.NewGameService()

	// Create template with functions
	templates := views.ParseTemplates("templates/game.html")

	// Set up routes
	http.HandleFunc("/", handlers.NewHandleIndex(gameService, templates))
	http.HandleFunc("POST /new-game", handlers.NewHandleNewGame(gameService))
	http.HandleFunc("POST /review-game", handlers.NewHandleReviewGame())
	http.HandleFunc("POST /select-card", handlers.NewHandleSelectCard(gameService))
	http.HandleFunc("POST /play-card", handlers.NewHandlePlayCard(gameService))
	http.HandleFunc("POST /ai-turn", handlers.NewHandleAITurn(gameService))

	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Start the server
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
