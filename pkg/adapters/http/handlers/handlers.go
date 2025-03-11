package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/xpmatteo/scopa-trainer/pkg/application"
	"github.com/xpmatteo/scopa-trainer/pkg/domain"
)

// Handler handles HTTP requests for the game
type Handler struct {
	service  *application.GameService
	template *template.Template
}

// NewHandler creates a new HTTP handler
func NewHandler(service *application.GameService, templ *template.Template) (*Handler, error) {
	return &Handler{
		service:  service,
		template: templ,
	}, nil
}

// HandleIndex serves the main game page
func (h *Handler) HandleIndex(w http.ResponseWriter, r *http.Request) {
	model := h.service.GetUIModel()
	if err := h.template.Execute(w, model); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// HandleNewGame handles the request to start a new game
func (h *Handler) HandleNewGame(w http.ResponseWriter, r *http.Request) {
	h.service.StartNewGame()

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// HandleSelectCard handles the selection of a card from the player's hand
func (h *Handler) HandleSelectCard(w http.ResponseWriter, r *http.Request) {
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
	h.service.SelectCard(suit, rank)

	// Get the current UI model to check if player's turn is complete
	model := h.service.GetUIModel()

	// If it's not the player's turn anymore (i.e., they made a capture),
	// trigger the AI turn
	if !model.PlayerTurn {
		h.service.PlayAITurn()
	}

	// Redirect to the main page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

type SelectedCardPlayer interface {
	PlaySelectedCard()
}

func NewHandlePlayCard(p SelectedCardPlayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p.PlaySelectedCard()

		// Redirect to the main page
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
