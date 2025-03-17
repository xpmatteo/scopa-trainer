package handlers

import (
	"html/template"
	"net/http"

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
