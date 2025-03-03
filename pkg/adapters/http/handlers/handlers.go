package handlers

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/xpmatteo/scopa-trainer/pkg/application"
)

// Handler handles HTTP requests for the game
type Handler struct {
	service  *application.GameService
	template *template.Template
}

// NewHandler creates a new HTTP handler
func NewHandler(service *application.GameService) (*Handler, error) {
	funcMap := template.FuncMap{
		"lower": strings.ToLower,
	}

	tmpl, err := template.New("game.html").Funcs(funcMap).ParseFiles("templates/game.html")
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

// HandleNewGame handles the request to start a new game
func (h *Handler) HandleNewGame(w http.ResponseWriter, r *http.Request) {
	model := h.service.StartNewGame()
	if err := h.template.Execute(w, model); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
