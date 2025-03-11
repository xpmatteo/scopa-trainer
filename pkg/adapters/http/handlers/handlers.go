package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/xpmatteo/scopa-trainer/pkg/domain"
)

// Handler handles HTTP requests for the game
type Handler struct {
	template *template.Template
}

// NewHandler creates a new HTTP handler
func NewHandler(templ *template.Template) (*Handler, error) {
	return &Handler{
		template: templ,
	}, nil
}

// UIModelProvider defines the interface for getting the UI model
type UIModelProvider interface {
	GetUIModel() domain.UIModel
}

// HandleIndex serves the main game page
func (h *Handler) HandleIndex(provider UIModelProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		model := provider.GetUIModel()
		if err := h.template.Execute(w, model); err != nil {
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
		// Only accept POST requests
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

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
