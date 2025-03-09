package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"html/template"

	"github.com/stretchr/testify/assert"
	"github.com/xpmatteo/scopa-trainer/pkg/application"
	"github.com/xpmatteo/scopa-trainer/pkg/domain"
)

func TestHandleNewGameRedirects(t *testing.T) {
	service := application.NewGameService()
	handler, err := NewHandler(service, nil)
	require.NoError(t, err)

	// Create a test request
	req := httptest.NewRequest("POST", "/new-game", nil)
	w := httptest.NewRecorder()

	// Call the handler
	handler.HandleNewGame(w, req)

	// Check the response
	resp := w.Result()
	defer resp.Body.Close()

	// Verify that we get a redirect response
	assert.Equal(t, http.StatusSeeOther, resp.StatusCode)

	// Verify the redirect location
	location, err := resp.Location()
	assert.NoError(t, err)
	assert.Equal(t, "/", location.String())
}

// TestHandleSelectCard tests the card selection handler
func TestHandleSelectCard(t *testing.T) {
	// Create service and start a game
	service := application.NewGameService()
	service.StartNewGame()

	// Create template for testing
	tmpl, err := template.ParseFiles("../../../../templates/game.html")
	require.NoError(t, err)

	// Create handler
	handler, err := NewHandler(service, tmpl)
	require.NoError(t, err)

	// Create a test request with suit and rank parameters
	req := httptest.NewRequest("GET", "/select-card?suit=Denari&rank=1", nil)
	w := httptest.NewRecorder()

	// Call the handler
	handler.HandleSelectCard(w, req)

	// Check the response
	resp := w.Result()
	defer resp.Body.Close()

	// Verify successful response
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Verify that the card was selected in the service
	model := service.GetUIModel()
	assert.NotEqual(t, domain.NO_CARD_SELECTED, model.SelectedCard)
}
