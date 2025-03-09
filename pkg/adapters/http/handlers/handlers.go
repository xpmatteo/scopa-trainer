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
	suit := domain.Suit(r.URL.Query().Get("suit"))
	rankStr := r.URL.Query().Get("rank")
	rankInt, err := strconv.Atoi(rankStr)
	if err != nil {
		http.Error(w, "Invalid rank parameter", http.StatusBadRequest)
		return
	}
	rank := domain.Rank(rankInt)

	model := h.service.SelectCard(suit, rank)

	if err := h.template.Execute(w, model); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
