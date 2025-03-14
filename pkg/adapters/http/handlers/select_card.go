package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/xpmatteo/scopa-trainer/pkg/domain"
)

// UIModelProvider defines the interface for getting the UI model
type UIModelProvider interface {
	GetUIModel() domain.UIModel
}

// NewHandleIndex creates a handler for the index page
func NewHandleIndex(provider UIModelProvider, templ *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		model := provider.GetUIModel()
		if err := templ.Execute(w, model); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

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
