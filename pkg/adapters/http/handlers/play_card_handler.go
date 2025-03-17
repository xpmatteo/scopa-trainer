package handlers

import (
	"net/http"
)

// SelectedCardPlayer defines the interface for playing a selected card
type SelectedCardPlayer interface {
	PlaySelectedCard()
}

// NewHandlePlayCard creates a handler for playing a selected card
func NewHandlePlayCard(player SelectedCardPlayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		player.PlaySelectedCard()

		// Redirect to the main page
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
