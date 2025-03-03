package handlers

import (
	"github.com/xpmatteo/scopa-trainer/pkg/application"
	"html/template"
	"net/http"
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
