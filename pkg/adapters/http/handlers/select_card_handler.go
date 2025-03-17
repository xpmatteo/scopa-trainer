package handlers

import (
	"net/http"
	"strconv"

	"github.com/xpmatteo/scopa-trainer/pkg/domain"
)

// CardSelector defines the interface for selecting a card
type CardSelector interface {
	SelectCard(suit domain.Suit, rank domain.Rank)
}

// NewHandleSelectCard creates a handler for selecting a card
func NewHandleSelectCard(selector CardSelector) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse form values
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		suit := domain.Suit(r.PostForm.Get("suit"))
		rankStr := r.PostForm.Get("rank")
		rankInt, err := strconv.Atoi(rankStr)
		if err != nil {
			http.Error(w, "Invalid rank parameter", http.StatusBadRequest)
			return
		}
		rank := domain.Rank(rankInt)

		// Process the action
		selector.SelectCard(suit, rank)

		// Redirect to the main page
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
