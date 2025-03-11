package handlers

import (
	"html/template"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xpmatteo/scopa-trainer/pkg/domain"
)

// FakeGameStarter is a fake implementation of the GameStarter interface
type FakeGameStarter struct {
	callCount int
}

func (f *FakeGameStarter) StartNewGame() {
	f.callCount++
}

// TestHandleNewGame tests the new game handler
func TestHandleNewGame(t *testing.T) {
	// Arrange
	fakeStarter := &FakeGameStarter{}
	handler := NewHandleNewGame(fakeStarter)

	// Create a test request
	req := httptest.NewRequest("POST", "/new-game", nil)
	w := httptest.NewRecorder()

	// Act
	handler.ServeHTTP(w, req)

	// Assert
	resp := w.Result()
	defer resp.Body.Close()

	// Verify that we get a redirect response
	assert.Equal(t, http.StatusSeeOther, resp.StatusCode)

	// Verify the redirect location
	location, err := resp.Location()
	assert.NoError(t, err)
	assert.Equal(t, "/", location.String())

	// Verify that the service method was called
	assert.Equal(t, 1, fakeStarter.callCount)
}

// FakeUIModelProvider is a fake implementation of the UIModelProvider interface
type FakeUIModelProvider struct {
	model domain.UIModel
}

func (f *FakeUIModelProvider) GetUIModel() domain.UIModel {
	return f.model
}

// TestHandleIndex tests the index handler
func TestHandleIndex(t *testing.T) {
	// Arrange
	// Create a simple template for testing
	tmpl, err := template.New("test").Parse("Game in progress: {{.GameInProgress}}")
	require.NoError(t, err)

	handler, err := NewHandler(tmpl)
	require.NoError(t, err)

	// Create a fake model provider
	model := domain.NewUIModel()
	model.GameInProgress = true
	provider := &FakeUIModelProvider{model: model}

	// Create a test request
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	// Act
	indexHandler := handler.HandleIndex(provider)
	indexHandler.ServeHTTP(w, req)

	// Assert
	resp := w.Result()
	defer resp.Body.Close()

	// Verify that we get a successful response
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Verify the response body contains the expected content
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.True(t, strings.Contains(string(body), "Game in progress: true"))
}

// FakeCardSelector is a fake implementation of the CardSelector interface
type FakeCardSelector struct {
	callCount int
	lastSuit  domain.Suit
	lastRank  domain.Rank
}

func (f *FakeCardSelector) SelectCard(suit domain.Suit, rank domain.Rank) {
	f.callCount++
	f.lastSuit = suit
	f.lastRank = rank
}

// TestHandleSelectCard tests the card selection handler
func TestHandleSelectCard(t *testing.T) {
	// Arrange
	fakeSelector := &FakeCardSelector{}
	handler := NewHandleSelectCard(fakeSelector)

	// Create a test request with POST method and form values
	req := httptest.NewRequest(http.MethodPost, "/select-card?suit=Coppe&rank=7", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.ParseForm()
	req.PostForm.Set("suit", "Coppe")
	req.PostForm.Set("rank", "7")
	w := httptest.NewRecorder()

	// Act
	handler.ServeHTTP(w, req)

	// Assert
	resp := w.Result()
	defer resp.Body.Close()

	// Verify we get a redirect response
	assert.Equal(t, http.StatusSeeOther, resp.StatusCode)

	// Verify the redirect location
	location, err := resp.Location()
	assert.NoError(t, err)
	assert.Equal(t, "/", location.String())

	// Verify that the service method was called with the correct parameters
	assert.Equal(t, 1, fakeSelector.callCount)
	assert.Equal(t, domain.Suit("Coppe"), fakeSelector.lastSuit)
	assert.Equal(t, domain.Rank(7), fakeSelector.lastRank)
}
