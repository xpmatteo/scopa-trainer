package handlers

import (
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"

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
