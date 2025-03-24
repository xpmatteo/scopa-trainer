package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xpmatteo/scopa-trainer/pkg/domain"
)

func TestCaptureCard(t *testing.T) {
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

	// Clear the table first - there might be initial cards there
	tableCardsInitial := service.gameState.Deck.CardsAt(domain.TableLocation)
	for _, card := range tableCardsInitial {
		service.gameState.Deck.MoveCard(card, domain.TableLocation, domain.DeckLocation)
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

	// Verify initial state
	initialModel := service.GetUIModel()
	assert.Equal(t, 1, len(initialModel.TableCards), "Table should have one card before capture")
	assert.Equal(t, 10, len(initialModel.PlayerHand), "Player should have 10 cards before capture")
	assert.Equal(t, 0, len(service.gameState.Deck.CardsAt(domain.PlayerCapturesLocation)), "Player should have no captured cards initially")

	// When the player selects a matching card from the table
	service.SelectCard(tableCard.Suit, tableCard.Rank)

	// Get the updated model
	model := service.GetUIModel()

	// Then the cards should be captured
	assert.Equal(t, 0, len(model.TableCards), "Table should have no cards after capture")
	assert.Equal(t, 9, len(model.PlayerHand), "Player should have 9 cards after capture")

	// And the cards should be in the player's capture pile
	capturedCards := service.gameState.Deck.CardsAt(domain.PlayerCapturesLocation)
	assert.Equal(t, 2, len(capturedCards), "Player should have 2 cards in capture pile")

	// And no card should be selected
	assert.Equal(t, domain.NO_CARD_SELECTED, model.SelectedCard, "No card should be selected after capture")

	// And it should be the AI's turn
	assert.Equal(t, domain.StatusAITurn, service.gameState.Status, "It should be the AI's turn after capture")

	// And player should have been awarded a scopa (cleared the table)
	assert.Equal(t, 1, service.playerScopaCount, "Player should have been awarded a scopa")
}

func TestCannotCaptureNonMatchingCard(t *testing.T) {
	// Given a game in progress with a card selected from hand
	service := NewGameService()
	service.StartNewGame()

	// Get a card from the player's hand
	playerHand := service.gameState.Deck.CardsAt(domain.PlayerHandLocation)
	selectedCard := playerHand[0]

	// Create a card with a different rank to put on the table
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

	// Get the initial state
	initialModel := service.GetUIModel()
	initialTableCardCount := len(initialModel.TableCards)

	// Manually set the selected card
	service.selectedCard = selectedCard

	// When the player selects a non-matching card from the table
	service.SelectCard(tableCard.Suit, tableCard.Rank)

	// Get the updated model
	model := service.GetUIModel()

	// Then the cards should not be captured
	assert.Equal(t, initialTableCardCount, len(model.TableCards), "Table should still have the same number of cards")
	assert.Equal(t, 10, len(model.PlayerHand), "Player should still have all cards")

	// And the selected card should remain selected
	assert.Equal(t, selectedCard, model.SelectedCard, "The hand card should remain selected")

	// And no cards should be in the capture pile
	capturedCards := service.gameState.Deck.CardsAt(domain.PlayerCapturesLocation)
	assert.Equal(t, 0, len(capturedCards), "Player should have no cards in capture pile")
}

func TestSelectingTableCardWithoutHandCardDoesNothing(t *testing.T) {
	// Given a game in progress with no card selected from hand
	service := NewGameService()
	service.StartNewGame()

	// Ensure no card is selected
	service.selectedCard = domain.NO_CARD_SELECTED

	// Put a card on the table
	tableCard := domain.Card{
		Suit: domain.Coppe,
		Rank: domain.Asso,
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
	initialTableCardCount := len(initialModel.TableCards)

	// Verify no card is selected
	assert.Equal(t, domain.NO_CARD_SELECTED, service.selectedCard, "No card should be selected initially")

	// When the player selects a card from the table without selecting a hand card first
	service.SelectCard(tableCard.Suit, tableCard.Rank)

	// Get the updated model
	model := service.GetUIModel()

	// Then nothing should happen
	assert.Equal(t, initialTableCardCount, len(model.TableCards), "Table should still have the same number of cards")
	assert.Equal(t, domain.NO_CARD_SELECTED, model.SelectedCard, "No card should be selected")
}
