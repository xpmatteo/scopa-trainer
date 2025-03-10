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

	// Get a card from the player's hand
	playerHand := service.gameState.Deck.CardsAt(domain.PlayerHandLocation)
	selectedCard := playerHand[0]

	// Create a card with the same rank to put on the table (making capture possible)
	tableCard := domain.Card{
		Suit:  domain.Coppe,
		Rank:  selectedCard.Rank,
		Name:  selectedCard.Rank.String(),
		Value: selectedCard.Rank.Value(),
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

	// Then the model should indicate that the card cannot be played to the table
	assert.False(t, model.CanPlaySelectedCard, "Should not be able to play card when capture is possible")

	// And when we try to play the card to the table
	beforeTableCount := len(model.TableCards)
	beforeHandCount := len(model.PlayerHand)

	playModel := service.PlaySelectedCard()

	// Then the card should not be played
	assert.Equal(t, beforeTableCount, len(playModel.TableCards), "Table card count should not change")
	assert.Equal(t, beforeHandCount, len(playModel.PlayerHand), "Hand card count should not change")
	assert.Equal(t, selectedCard, playModel.SelectedCard, "Selected card should remain selected")
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
		Suit:  domain.Coppe,
		Rank:  differentRank,
		Name:  differentRank.String(),
		Value: differentRank.Value(),
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

	// And when we try to play the card to the table
	beforeTableCount := len(model.TableCards)
	beforeHandCount := len(model.PlayerHand)

	playModel := service.PlaySelectedCard()

	// Then the card should be played
	assert.Equal(t, beforeTableCount+1, len(playModel.TableCards), "Table should have one more card")
	assert.Equal(t, beforeHandCount-1, len(playModel.PlayerHand), "Hand should have one less card")
	assert.Equal(t, domain.NO_CARD_SELECTED, playModel.SelectedCard, "No card should be selected after playing")
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

	// And when we try to play the card to the table
	beforeTableCount := len(model.TableCards)
	beforeHandCount := len(model.PlayerHand)

	playModel := service.PlaySelectedCard()

	// Then the card should be played
	assert.Equal(t, beforeTableCount+1, len(playModel.TableCards), "Table should have one more card")
	assert.Equal(t, beforeHandCount-1, len(playModel.PlayerHand), "Hand should have one less card")
	assert.Equal(t, domain.NO_CARD_SELECTED, playModel.SelectedCard, "No card should be selected after playing")
}
