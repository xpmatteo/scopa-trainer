package handlers

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xpmatteo/scopa-trainer/pkg/application"
)

func TestHandlePlayCard(t *testing.T) {
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

	// Select a card from the player's hand first
	playerHand := service.GetUIModel().PlayerHand
	require.NotEmpty(t, playerHand, "Player hand should not be empty")
	selectedCard := playerHand[0]
	service.SelectCard(selectedCard.Suit, selectedCard.Rank)

	// Create a test request for playing the card
	req := httptest.NewRequest(http.MethodPost, "/play-card", nil)
	w := httptest.NewRecorder()

	// Call the handler
	handler.HandlePlayCard(w, req)

	// Check the response
	resp := w.Result()
	defer resp.Body.Close()

	// Verify we get a redirect response
	assert.Equal(t, http.StatusSeeOther, resp.StatusCode)

	// Verify the redirect location
	location, err := resp.Location()
	assert.NoError(t, err)
	assert.Equal(t, "/", location.String())

	// Verify that the card was played to the table
	model := service.GetUIModel()
	assert.Equal(t, 1, len(model.TableCards), "Table should have one card after playing")
	assert.Equal(t, 9, len(model.PlayerHand), "Player should have 9 cards after playing one")
}
