package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xpmatteo/scopa-trainer/pkg/domain"
)

func TestCardCountsInUIModel(t *testing.T) {
	// Given a new game
	service := NewGameService()
	service.StartNewGame()

	// When we get the UI model
	model := service.GetUIModel()

	// Then the card counts should be correct
	// 40 total cards - 10 player hand - 10 AI hand = 20 in deck
	assert.Equal(t, 20, model.DeckCount, "Deck should have 20 cards at start")
	assert.Equal(t, 0, model.PlayerCaptureCount, "Player should have 0 captured cards at start")
	assert.Equal(t, 0, model.AICaptureCount, "AI should have 0 captured cards at start")
}

func TestCardCountsAfterCapture(t *testing.T) {
	// Given a game in progress with a card selected from hand
	service := NewGameService()
	service.StartNewGame()

	// Get a card from the player's hand
	playerHand := service.gameState.Deck.CardsAt(domain.PlayerHandLocation)
	selectedCard := playerHand[0]

	// Create a card with the same rank to put on the table
	tableCard := domain.Card{
		Suit: domain.Coppe,
		Rank: selectedCard.Rank,
	}

	// Find this card in the deck and move it to the table
	found := false
	for _, card := range service.gameState.Deck.CardsAt(domain.DeckLocation) {
		if card.Rank == tableCard.Rank && card.Suit == tableCard.Suit {
			service.gameState.Deck.MoveCard(card, domain.DeckLocation, domain.TableLocation)
			tableCard = card // Use the actual card from the deck
			found = true
			break
		}
	}
	if !found {
		// Try to find any card with the same rank
		for _, card := range service.gameState.Deck.CardsAt(domain.DeckLocation) {
			if card.Rank == tableCard.Rank {
				service.gameState.Deck.MoveCard(card, domain.DeckLocation, domain.TableLocation)
				tableCard = card // Use the actual card from the deck
				break
			}
		}
	}

	// Get the initial state
	initialModel := service.GetUIModel()
	initialDeckCount := initialModel.DeckCount

	// Manually set the selected card
	service.selectedCard = selectedCard

	// When the player captures a card
	service.SelectCard(tableCard.Suit, tableCard.Rank)

	// Then the card counts should be updated
	model := service.GetUIModel()
	assert.Equal(t, initialDeckCount, model.DeckCount, "Deck count should not change after capture")
	assert.Equal(t, 2, model.PlayerCaptureCount, "Player should have 2 captured cards")
	assert.Equal(t, 0, model.AICaptureCount, "AI should still have 0 captured cards")
}

func TestCardCountsAfterPlayingToTable(t *testing.T) {
	// Given a game in progress with a card selected from hand
	service := NewGameService()
	service.StartNewGame()

	// Get a card from the player's hand
	playerHand := service.gameState.Deck.CardsAt(domain.PlayerHandLocation)
	selectedCard := playerHand[0]

	// Get the initial state
	initialModel := service.GetUIModel()
	initialDeckCount := initialModel.DeckCount

	// Manually set the selected card
	service.selectedCard = selectedCard

	// When the player plays a card to the table
	service.PlaySelectedCard()

	// Then the card counts should be updated
	model := service.GetUIModel()
	assert.Equal(t, initialDeckCount, model.DeckCount, "Deck count should not change after playing to table")
	assert.Equal(t, 0, model.PlayerCaptureCount, "Player should still have 0 captured cards")
	assert.Equal(t, 0, model.AICaptureCount, "AI should still have 0 captured cards")
}
