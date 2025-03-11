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

	// Create a fake model provider
	model := domain.NewUIModel()
	model.GameInProgress = true
	provider := &FakeUIModelProvider{model: model}

	// Create handler
	handler := NewHandleIndex(provider, tmpl)

	// Create a test request
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	// Act
	handler.ServeHTTP(w, req)

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
