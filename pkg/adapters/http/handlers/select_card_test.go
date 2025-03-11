package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xpmatteo/scopa-trainer/pkg/domain"
)

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
