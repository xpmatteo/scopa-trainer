package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"html/template"

	"net/url"
	"strconv"
	"strings"

	"github.com/stretchr/testify/assert"
	"github.com/xpmatteo/scopa-trainer/pkg/application"
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
	funcMap := template.FuncMap{
		"lower": strings.ToLower,
	}
	tmpl, err := template.New("game.html").Funcs(funcMap).ParseFiles("../../../../templates/game.html")
	require.NoError(t, err)

	// Create handler
	handler, err := NewHandler(service, tmpl)
	require.NoError(t, err)

	// Get a card from the player's hand to use in the test
	playerHand := service.GetUIModel().PlayerHand
	require.NotEmpty(t, playerHand, "Player hand should not be empty")
	selectedCard := playerHand[0]

	// Create a test request with POST method and form values
	form := url.Values{}
	form.Add("suit", string(selectedCard.Suit))
	form.Add("rank", strconv.Itoa(int(selectedCard.Rank)))
	req := httptest.NewRequest(http.MethodPost, "/select-card", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	// Call the handler
	handler.HandleSelectCard(w, req)

	// Check the response
	resp := w.Result()
	defer resp.Body.Close()

	// Verify we get a redirect response
	assert.Equal(t, http.StatusSeeOther, resp.StatusCode)

	// Verify the redirect location
	location, err := resp.Location()
	assert.NoError(t, err)
	assert.Equal(t, "/", location.String())

	// Verify that the card was selected in the service
	model := service.GetUIModel()
	assert.Equal(t, selectedCard, model.SelectedCard, "The selected card should match the one we chose")
}
