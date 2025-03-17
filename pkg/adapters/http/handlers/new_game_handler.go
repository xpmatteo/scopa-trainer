package handlers

import (
	"net/http"
)

// GameStarter defines the interface for starting a new game
type GameStarter interface {
	StartNewGame()
}

// NewHandleNewGame creates a handler for starting a new game
func NewHandleNewGame(starter GameStarter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		starter.StartNewGame()
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
