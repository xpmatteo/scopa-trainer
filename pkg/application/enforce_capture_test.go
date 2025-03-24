package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xpmatteo/scopa-trainer/pkg/domain"
)

func TestCannotPlayCardWhenCaptureIsPossible(t *testing.T) {
	// Given a game in progress with a card selected from hand
	service := NewGameService()
	service.StartNewGame()

	// Use our test combination setup for a clean test environment
	service.SetupCombinationTest()

	// Add a specific card to the table that can be captured
	service.gameState.Deck.MoveCardMatching(domain.DeckLocation, domain.TableLocation, domain.Quattro, domain.Coppe)

	// Add a matching card to the player's hand
	service.gameState.Deck.MoveCardMatching(domain.DeckLocation, domain.PlayerHandLocation, domain.Quattro, domain.Spade)

	// Get the player card for our test
	playerHandCards := service.gameState.Deck.CardsAt(domain.PlayerHandLocation)
	if len(playerHandCards) != 1 {
		t.Fatalf("Expected 1 card in player's hand, got %d", len(playerHandCards))
	}
	playerCard := playerHandCards[0]

	// Verify initial state
	tableCardCount := len(service.gameState.Deck.CardsAt(domain.TableLocation))
	playerHandCount := len(service.gameState.Deck.CardsAt(domain.PlayerHandLocation))

	t.Logf("Initial state: table=%d, playerHand=%d", tableCardCount, playerHandCount)
	assert.Equal(t, 1, tableCardCount, "Table should have one card")
	assert.Equal(t, 1, playerHandCount, "Player should have one card")

	// Make sure it's player's turn
	service.gameState.Status = domain.StatusPlayerTurn

	// Select the card from hand
	service.selectedCard = playerCard

	// When the player tries to play the card to the table instead of capturing
	canPlay := !service.canCaptureAnyCard(service.selectedCard)
	assert.False(t, canPlay, "Should not be able to play card when capture is possible")

	// Play card despite capture being possible (directly testing the PlaySelectedCard method)
	service.PlaySelectedCard()

	// Then the card should not be played (nothing should change)
	postTableCardCount := len(service.gameState.Deck.CardsAt(domain.TableLocation))
	postPlayerHandCount := len(service.gameState.Deck.CardsAt(domain.PlayerHandLocation))

	t.Logf("After attempted play: table=%d, playerHand=%d, selectedCard=%v, status=%v",
		postTableCardCount, postPlayerHandCount, service.selectedCard, service.gameState.Status)

	assert.Equal(t, tableCardCount, postTableCardCount, "Table card count should not change")
	assert.Equal(t, playerHandCount, postPlayerHandCount, "Hand card count should not change")
	assert.Equal(t, playerCard, service.selectedCard, "Selected card should remain selected")
	assert.Equal(t, domain.StatusPlayerTurn, service.gameState.Status, "Game status should not change")
}

func TestCanPlayCardWhenNoCaptureIsPossible(t *testing.T) {
	// Given a game in progress with a card selected from hand
	service := NewGameService()
	service.StartNewGame()

	// Get a card from the player's hand
	playerHand := service.gameState.Deck.CardsAt(domain.PlayerHandLocation)
	selectedCard := playerHand[0]

	// Create a card with a different rank to put on the table (making capture impossible)
	differentRank := selectedCard.Rank + 1
	if differentRank > domain.Re {
		differentRank = domain.Asso
	}

	tableCard := domain.Card{
		Suit: domain.Coppe,
		Rank: differentRank,
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

	// Manually set the selected card
	service.selectedCard = selectedCard

	// When we get the UI model
	model := service.GetUIModel()

	// Then the model should indicate that the card can be played to the table
	assert.True(t, model.CanPlaySelectedCard, "Should be able to play card when no capture is possible")

	// Get the state before trying to play
	beforeTableCount := len(model.TableCards)
	beforeHandCount := len(model.PlayerHand)

	// And when we try to play the card to the table
	service.PlaySelectedCard()

	// Get the updated model
	playModel := service.GetUIModel()

	// Then the card should be played
	assert.Equal(t, beforeTableCount+1, len(playModel.TableCards), "Table should have one more card")
	assert.Equal(t, beforeHandCount-1, len(playModel.PlayerHand), "Hand should have one less card")
	assert.Equal(t, domain.NO_CARD_SELECTED, playModel.SelectedCard, "No card should be selected after playing")
	assert.Equal(t, domain.StatusAITurn, service.gameState.Status, "It should be AI's turn after playing a card")
}

func TestCanPlayCardWhenTableIsEmpty(t *testing.T) {
	// Given a game in progress with a card selected from hand and empty table
	service := NewGameService()
	service.StartNewGame()

	// Get a card from the player's hand
	playerHand := service.gameState.Deck.CardsAt(domain.PlayerHandLocation)
	selectedCard := playerHand[0]

	// Manually set the selected card
	service.selectedCard = selectedCard

	// When we get the UI model
	model := service.GetUIModel()

	// Then the model should indicate that the card can be played to the table
	assert.True(t, model.CanPlaySelectedCard, "Should be able to play card when table is empty")

	// Get the state before trying to play
	beforeTableCount := len(model.TableCards)
	beforeHandCount := len(model.PlayerHand)

	// And when we try to play the card to the table
	service.PlaySelectedCard()

	// Get the updated model
	playModel := service.GetUIModel()

	// Then the card should be played
	assert.Equal(t, beforeTableCount+1, len(playModel.TableCards), "Table should have one more card")
	assert.Equal(t, beforeHandCount-1, len(playModel.PlayerHand), "Hand should have one less card")
	assert.Equal(t, domain.NO_CARD_SELECTED, playModel.SelectedCard, "No card should be selected after playing")
	assert.Equal(t, domain.StatusAITurn, service.gameState.Status, "It should be AI's turn after playing a card")
}
