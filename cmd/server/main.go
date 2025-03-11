package main

import (
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/xpmatteo/scopa-trainer/pkg/adapters/http/handlers"
	"github.com/xpmatteo/scopa-trainer/pkg/application"
)

func main() {
	// Initialize the application service
	gameService := application.NewGameService()

	funcMap := template.FuncMap{
		"lower": strings.ToLower,
	}
	templ, err := template.New("game.html").Funcs(funcMap).ParseFiles("templates/game.html")
	if err != nil {
		panic(err)
	}

	// Set up routes
	http.HandleFunc("GET /", handlers.NewHandleIndex(gameService, templ))
	http.HandleFunc("POST /new-game", handlers.NewHandleNewGame(gameService))
	http.HandleFunc("POST /select-card", handlers.NewHandleSelectCard(gameService))
	http.HandleFunc("POST /play-card", handlers.NewHandlePlayCard(gameService))

	// Start the server
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
