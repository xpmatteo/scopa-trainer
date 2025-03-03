package handlers

import (
	"html/template"
	"net/http"

	"github.com/xpmatteo/scopa-trainer/pkg/application"
)

// Handler handles HTTP requests for the game
type Handler struct {
	service  *application.GameService
	template *template.Template
}

// NewHandler creates a new HTTP handler
func NewHandler(service *application.GameService) (*Handler, error) {
	tmpl, err := template.ParseFiles("templates/game.html")
	if err != nil {
		return nil, err
	}

	return &Handler{
		service:  service,
		template: tmpl,
	}, nil
}

// HandleIndex serves the main game page
func (h *Handler) HandleIndex(w http.ResponseWriter, r *http.Request) {
	model := h.service.GetUIModel()
	if err := h.template.Execute(w, model); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
