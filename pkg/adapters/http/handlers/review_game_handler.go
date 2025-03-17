package handlers

import (
	"net/http"
)

// NewHandleReviewGame creates a handler for reviewing the game
// This is a placeholder for future functionality
func NewHandleReviewGame() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// For now, just redirect back to the game page
		// This is a placeholder for future functionality
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
