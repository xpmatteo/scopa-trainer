package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xpmatteo/scopa-trainer/pkg/domain"
)

func TestPlayCardToTable(t *testing.T) {
	// Given a game in progress with a card selected from hand
	service := NewGameService()
	service.StartNewGame()

	// Get a card from the player's hand
	playerHand := service.gameState.Deck.CardsAt(domain.PlayerHandLocation)
	selectedCard := playerHand[0]

	// Select the card from hand
	service.selectedCard = selectedCard

	// When the player plays the card to the table
	model := service.PlaySelectedCard()

	// Then the card should be moved to the table
	assert.Equal(t, 1, len(model.TableCards), "Table should have one card after playing")
	assert.Equal(t, 9, len(model.PlayerHand), "Player should have 9 cards after playing one")

	// And the card on the table should be the one that was in the hand
	assert.Equal(t, selectedCard, model.TableCards[0], "The card on the table should be the one from the hand")

	// And no card should be selected
	assert.Equal(t, domain.NO_CARD_SELECTED, model.SelectedCard, "No card should be selected after playing")
}

func TestCannotPlayCardIfNoneSelected(t *testing.T) {
	// Given a game in progress with no card selected
	service := NewGameService()
	service.StartNewGame()
	service.selectedCard = domain.NO_CARD_SELECTED

	// When the player tries to play a card
	model := service.PlaySelectedCard()

	// Then no card should be played
	assert.Equal(t, 0, len(model.TableCards), "Table should still have no cards")
	assert.Equal(t, 10, len(model.PlayerHand), "Player should still have all 10 cards")

	// And no card should be selected
	assert.Equal(t, domain.NO_CARD_SELECTED, model.SelectedCard, "No card should be selected")
}
