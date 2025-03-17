package handlers

import (
	"net/http"
)

// AITurnPlayer defines the interface for playing the AI's turn
type AITurnPlayer interface {
	PlayAITurn()
}

// NewHandleAITurn creates a handler for triggering the AI's turn
func NewHandleAITurn(player AITurnPlayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		player.PlayAITurn()

		// Redirect to the main page
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
