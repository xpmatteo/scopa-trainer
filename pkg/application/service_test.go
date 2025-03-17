package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xpmatteo/scopa-trainer/pkg/domain"
)

// TestSortedCaptureCards tests that cards in the capture piles are sorted in the UI model
func TestSortedCaptureCards(t *testing.T) {
	// Given a game service with some capture cards in random order
	service := NewGameService()
	service.StartNewGame()

	// Move some cards from deck to player captures in random order
	deckCards := service.gameState.Deck.CardsAt(domain.DeckLocation)

	// Ensure we have enough cards in the deck
	assert.GreaterOrEqual(t, len(deckCards), 10, "Not enough cards in deck for test")

	// Pick some cards with different ranks and suits - use fewer cards to avoid index out of range
	cardsToMove := []domain.Card{
		deckCards[0],
		deckCards[1],
		deckCards[2],
		deckCards[3],
		deckCards[4],
	}

	// Move cards to player captures
	for _, card := range cardsToMove {
		service.gameState.Deck.MoveCard(card, domain.DeckLocation, domain.PlayerCapturesLocation)
	}

	// Move some cards to AI captures as well
	for i := 5; i < 10 && i < len(deckCards); i++ {
		service.gameState.Deck.MoveCard(deckCards[i], domain.DeckLocation, domain.AICapturesLocation)
	}

	// When we get the UI model
	model := service.GetUIModel()

	// Then the capture cards should be sorted by rank and suit
	assert.NotEmpty(t, model.PlayerCaptureCards)
	assert.NotEmpty(t, model.AICaptureCards)

	// Verify player capture cards are sorted
	for i := 1; i < len(model.PlayerCaptureCards); i++ {
		// Either the current card's rank is greater than the previous card's rank,
		// or if ranks are equal, the current card's suit should come later alphabetically
		prevCard := model.PlayerCaptureCards[i-1]
		currCard := model.PlayerCaptureCards[i]

		if prevCard.Rank == currCard.Rank {
			assert.LessOrEqual(t, string(prevCard.Suit), string(currCard.Suit),
				"Cards with same rank should be sorted by suit")
		} else {
			assert.LessOrEqual(t, prevCard.Rank, currCard.Rank,
				"Cards should be sorted by rank")
		}
	}

	// Verify AI capture cards are sorted
	for i := 1; i < len(model.AICaptureCards); i++ {
		prevCard := model.AICaptureCards[i-1]
		currCard := model.AICaptureCards[i]

		if prevCard.Rank == currCard.Rank {
			assert.LessOrEqual(t, string(prevCard.Suit), string(currCard.Suit),
				"Cards with same rank should be sorted by suit")
		} else {
			assert.LessOrEqual(t, prevCard.Rank, currCard.Rank,
				"Cards should be sorted by rank")
		}
	}
}
