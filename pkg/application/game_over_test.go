package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xpmatteo/scopa-trainer/pkg/domain"
)

func TestGameOverWhenHandsAndDeckAreEmpty(t *testing.T) {
	// Given a game in progress with empty hands and an empty deck
	service := NewGameService()
	service.StartNewGame()

	// Move all cards from player and AI hands to captures (to empty the hands)
	playerHand := service.gameState.Deck.CardsAt(domain.PlayerHandLocation)
	for _, card := range playerHand {
		service.gameState.Deck.MoveCard(card, domain.PlayerHandLocation, domain.PlayerCapturesLocation)
	}

	aiHand := service.gameState.Deck.CardsAt(domain.AIHandLocation)
	for _, card := range aiHand {
		service.gameState.Deck.MoveCard(card, domain.AIHandLocation, domain.AICapturesLocation)
	}

	// Move all cards from deck to table (to empty the deck)
	deckCards := service.gameState.Deck.CardsAt(domain.DeckLocation)
	for _, card := range deckCards {
		service.gameState.Deck.MoveCard(card, domain.DeckLocation, domain.TableLocation)
	}

	// Verify hands and deck are empty
	assert.Equal(t, 0, len(service.gameState.Deck.CardsAt(domain.PlayerHandLocation)))
	assert.Equal(t, 0, len(service.gameState.Deck.CardsAt(domain.AIHandLocation)))
	assert.Equal(t, 0, len(service.gameState.Deck.CardsAt(domain.DeckLocation)))

	// When we check if new cards need to be dealt
	service.DealNewCardsIfNeeded()

	// Then the game should be over
	assert.Equal(t, domain.StatusGameOver, service.gameState.Status, "Game should be over when hands and deck are empty")

	// And the UI model should reflect the game over state
	model := service.GetUIModel()
	assert.True(t, model.GameOver, "UI model should indicate game is over")
	assert.True(t, model.ShowNewGameButton, "New Game button should be displayed when game is over")
}

func TestGameNotOverWhenOnlyHandsAreEmpty(t *testing.T) {
	// Given a game in progress with empty hands but cards in the deck
	service := NewGameService()
	service.StartNewGame()

	// Move all cards from player and AI hands to captures (to empty the hands)
	playerHand := service.gameState.Deck.CardsAt(domain.PlayerHandLocation)
	for _, card := range playerHand {
		service.gameState.Deck.MoveCard(card, domain.PlayerHandLocation, domain.PlayerCapturesLocation)
	}

	aiHand := service.gameState.Deck.CardsAt(domain.AIHandLocation)
	for _, card := range aiHand {
		service.gameState.Deck.MoveCard(card, domain.AIHandLocation, domain.AICapturesLocation)
	}

	// Verify hands are empty but deck has cards
	assert.Equal(t, 0, len(service.gameState.Deck.CardsAt(domain.PlayerHandLocation)))
	assert.Equal(t, 0, len(service.gameState.Deck.CardsAt(domain.AIHandLocation)))
	assert.Greater(t, len(service.gameState.Deck.CardsAt(domain.DeckLocation)), 0)

	// When we check if new cards need to be dealt
	service.DealNewCardsIfNeeded()

	// Then the game should not be over, and new cards should be dealt
	assert.NotEqual(t, domain.StatusGameOver, service.gameState.Status, "Game should not be over when deck has cards")

	// And the UI model should not reflect the game over state
	model := service.GetUIModel()
	assert.False(t, model.GameOver, "UI model should not indicate game is over")
}

func TestGameNotOverWhenOnlyDeckIsEmpty(t *testing.T) {
	// Given a game in progress with cards in hands but an empty deck
	service := NewGameService()
	service.StartNewGame()

	// Move all cards from deck to table (to empty the deck)
	deckCards := service.gameState.Deck.CardsAt(domain.DeckLocation)
	for _, card := range deckCards {
		service.gameState.Deck.MoveCard(card, domain.DeckLocation, domain.TableLocation)
	}

	// Verify deck is empty but hands have cards
	assert.Equal(t, 0, len(service.gameState.Deck.CardsAt(domain.DeckLocation)))
	assert.Greater(t, len(service.gameState.Deck.CardsAt(domain.PlayerHandLocation)), 0)
	assert.Greater(t, len(service.gameState.Deck.CardsAt(domain.AIHandLocation)), 0)

	// When we check if new cards need to be dealt
	service.DealNewCardsIfNeeded()

	// Then the game should not be over
	assert.NotEqual(t, domain.StatusGameOver, service.gameState.Status, "Game should not be over when hands have cards")

	// And the UI model should not reflect the game over state
	model := service.GetUIModel()
	assert.False(t, model.GameOver, "UI model should not indicate game is over")
}
