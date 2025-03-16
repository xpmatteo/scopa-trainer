package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xpmatteo/scopa-trainer/pkg/domain"
)

func TestDealNewCardsWhenHandsAreEmpty(t *testing.T) {
	// Given a game in progress with empty hands and cards in the deck
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

	// Verify hands are empty and deck has cards
	assert.Equal(t, 0, len(service.gameState.Deck.CardsAt(domain.PlayerHandLocation)))
	assert.Equal(t, 0, len(service.gameState.Deck.CardsAt(domain.AIHandLocation)))
	assert.Equal(t, 20, len(service.gameState.Deck.CardsAt(domain.DeckLocation)))

	// When we check if new cards need to be dealt
	cardsDealt := service.DealNewCardsIfNeeded()

	// Then new cards should be dealt to both players
	assert.True(t, cardsDealt, "Cards should have been dealt")
	assert.Equal(t, 10, len(service.gameState.Deck.CardsAt(domain.PlayerHandLocation)))
	assert.Equal(t, 10, len(service.gameState.Deck.CardsAt(domain.AIHandLocation)))
	assert.Equal(t, 0, len(service.gameState.Deck.CardsAt(domain.DeckLocation)))
}

func TestDealNewCardsWhenDeckIsEmpty(t *testing.T) {
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
	cardsDealt := service.DealNewCardsIfNeeded()

	// Then no cards should be dealt
	assert.False(t, cardsDealt, "No cards should have been dealt")
	assert.Equal(t, 0, len(service.gameState.Deck.CardsAt(domain.PlayerHandLocation)))
	assert.Equal(t, 0, len(service.gameState.Deck.CardsAt(domain.AIHandLocation)))
}

func TestDealNewCardsWhenOnlyPlayerHandIsEmpty(t *testing.T) {
	// Given a game in progress with only player's hand empty
	service := NewGameService()
	service.StartNewGame()

	// Move all cards from player hand to captures (to empty the hand)
	playerHand := service.gameState.Deck.CardsAt(domain.PlayerHandLocation)
	for _, card := range playerHand {
		service.gameState.Deck.MoveCard(card, domain.PlayerHandLocation, domain.PlayerCapturesLocation)
	}

	// Verify player hand is empty but AI hand is not
	assert.Equal(t, 0, len(service.gameState.Deck.CardsAt(domain.PlayerHandLocation)))
	assert.Equal(t, 10, len(service.gameState.Deck.CardsAt(domain.AIHandLocation)))

	// When we check if new cards need to be dealt
	cardsDealt := service.DealNewCardsIfNeeded()

	// Then no cards should be dealt (both hands must be empty)
	assert.False(t, cardsDealt, "No cards should have been dealt")
	assert.Equal(t, 0, len(service.gameState.Deck.CardsAt(domain.PlayerHandLocation)))
	assert.Equal(t, 10, len(service.gameState.Deck.CardsAt(domain.AIHandLocation)))
}

func TestDealNewCardsAfterPlayingLastCard(t *testing.T) {
	// Given a game in progress with only one card in each hand
	service := NewGameService()
	service.StartNewGame()

	// Move all but one card from player and AI hands to captures
	playerHand := service.gameState.Deck.CardsAt(domain.PlayerHandLocation)
	for i := 0; i < len(playerHand)-1; i++ {
		service.gameState.Deck.MoveCard(playerHand[i], domain.PlayerHandLocation, domain.PlayerCapturesLocation)
	}

	aiHand := service.gameState.Deck.CardsAt(domain.AIHandLocation)
	for i := 0; i < len(aiHand)-1; i++ {
		service.gameState.Deck.MoveCard(aiHand[i], domain.AIHandLocation, domain.AICapturesLocation)
	}

	// Verify each hand has only one card
	playerHand = service.gameState.Deck.CardsAt(domain.PlayerHandLocation)
	aiHand = service.gameState.Deck.CardsAt(domain.AIHandLocation)
	assert.Equal(t, 1, len(playerHand))
	assert.Equal(t, 1, len(aiHand))

	// Select the player's last card
	service.selectedCard = playerHand[0]

	// When the player plays their last card
	service.PlaySelectedCard()

	// Verify player's hand is empty and it's AI's turn
	assert.Equal(t, 0, len(service.gameState.Deck.CardsAt(domain.PlayerHandLocation)))
	assert.Equal(t, domain.AITurn, service.gameState.Status, "It should be AI's turn")

	// Then manually trigger the AI turn
	service.PlayAITurn()

	// Verify both hands are empty and new cards have been dealt
	assert.Equal(t, 10, len(service.gameState.Deck.CardsAt(domain.PlayerHandLocation)))
	assert.Equal(t, 10, len(service.gameState.Deck.CardsAt(domain.AIHandLocation)))
	assert.Equal(t, 0, len(service.gameState.Deck.CardsAt(domain.DeckLocation)))
}
