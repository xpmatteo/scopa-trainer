package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type FakeService struct {
	callCount int
}

func (f *FakeService) PlaySelectedCard() {
	f.callCount++
}

func TestHandlePlayCard(t *testing.T) {
	// Arrange
	// Create a new handler
	fakeService := &FakeService{}
	handler := NewHandlePlayCard(fakeService)

	// Create a test request for playing the card
	req := httptest.NewRequest(http.MethodPost, "/play-card", nil)
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
