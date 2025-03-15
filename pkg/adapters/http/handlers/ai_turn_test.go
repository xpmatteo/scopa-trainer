package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type FakeAIPlayer struct {
	callCount int
}

func (f *FakeAIPlayer) PlayAITurn() {
	f.callCount++
}

func TestHandleAITurn(t *testing.T) {
	// Arrange
	// Create a new handler
	fakeService := &FakeAIPlayer{}
	handler := NewHandleAITurn(fakeService)

	// Create a test request for playing the AI turn
	req := httptest.NewRequest(http.MethodPost, "/ai-turn", nil)
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

	// Verify that the service method was called
	assert.Equal(t, 1, fakeService.callCount)
}
