package handlers

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
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
