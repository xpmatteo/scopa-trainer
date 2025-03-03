package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDeckCreatesAllCardsInDeckLocation(t *testing.T) {
	deck := NewDeck()

	// A new deck should have all 40 cards in the DeckLocation
	assert.Equal(t, 40, len(deck.CardsAt(DeckLocation)))

	// No cards should be in other locations
	assert.Equal(t, 0, len(deck.CardsAt(PlayerHandLocation)))
	assert.Equal(t, 0, len(deck.CardsAt(AIHandLocation)))
	assert.Equal(t, 0, len(deck.CardsAt(TableLocation)))
	assert.Equal(t, 0, len(deck.CardsAt(PlayerCapturesLocation)))
	assert.Equal(t, 0, len(deck.CardsAt(AICapturesLocation)))
}

func TestDealCards(t *testing.T) {
	deck := NewDeck()

	// Deal 10 cards to the player
	deck.DealCards(DeckLocation, PlayerHandLocation, 10)

	// Check that 10 cards moved to player hand and 30 remain in deck
	assert.Equal(t, 30, len(deck.CardsAt(DeckLocation)))
	assert.Equal(t, 10, len(deck.CardsAt(PlayerHandLocation)))

	// Deal 10 cards to the AI
	deck.DealCards(DeckLocation, AIHandLocation, 10)

	// Check that 10 cards moved to AI hand and 20 remain in deck
	assert.Equal(t, 20, len(deck.CardsAt(DeckLocation)))
	assert.Equal(t, 10, len(deck.CardsAt(AIHandLocation)))
}

func TestMoveCard(t *testing.T) {
	deck := NewDeck()

	// Deal cards to player hand
	deck.DealCards(DeckLocation, PlayerHandLocation, 10)
	playerCards := deck.CardsAt(PlayerHandLocation)
	cardToMove := playerCards[0]

	// Move one card from player hand to table
	deck.MoveCard(cardToMove, PlayerHandLocation, TableLocation)

	// Check that the card moved from player hand to table
	assert.Equal(t, 9, len(deck.CardsAt(PlayerHandLocation)))
	assert.Equal(t, 1, len(deck.CardsAt(TableLocation)))
	assert.Equal(t, TableLocation, deck.GetCardLocation(cardToMove))
}

func TestShuffling(t *testing.T) {
	deck := NewDeck()

	// Get the original order
	originalOrder := deck.CardsAt(DeckLocation)

	// Shuffle the deck
	deck.Shuffle()

	// Get the new order
	newOrder := deck.CardsAt(DeckLocation)

	// Check that all cards are still in the deck
	assert.Equal(t, 40, len(newOrder))

	// Count differences to verify shuffle happened
	differences := 0
	for i := range originalOrder {
		if i < len(newOrder) && originalOrder[i] != newOrder[i] {
			differences++
		}
	}

	// There should be a substantial number of differences after shuffling
	assert.True(t, differences > 10, "Shuffle didn't change card order significantly")
}
