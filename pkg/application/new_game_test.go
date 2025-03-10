package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/xpmatteo/scopa-trainer/pkg/domain"
)

type NewGameTestSuite struct {
	suite.Suite
	service *GameService
}

func (s *NewGameTestSuite) SetupTest() {
	s.service = NewGameService()
}

func TestNewGameSuite(t *testing.T) {
	suite.Run(t, new(NewGameTestSuite))
}

func TestInitialUIModelHasNewGameButton(t *testing.T) {
	service := NewGameService()

	// When the application starts
	model := service.GetUIModel()

	// Then the UI model should have a "New Game" button
	assert.True(t, model.ShowNewGameButton, "New Game button should be displayed initially")
}

func TestStartNewGame(t *testing.T) {
	// Given the application has started
	service := NewGameService()
	initialModel := service.GetUIModel()
	assert.True(t, initialModel.ShowNewGameButton, "New Game button should be displayed initially")

	// When the player clicks the "New Game" button
	service.StartNewGame()

	// Get the updated model
	updatedModel := service.GetUIModel()

	assert.Equal(t, 0, len(updatedModel.TableCards), "Table should have no cards")
	assert.Equal(t, 10, len(updatedModel.PlayerHand), "Player should have 10 cards")
	assert.True(t, isSorted(updatedModel.PlayerHand), "Player hand should be sorted")
	assert.False(t, updatedModel.ShowNewGameButton, "New Game button should be hidden after starting a game")
}

func isSorted(cards []domain.Card) bool {
	for i := 1; i < len(cards); i++ {
		if cards[i].Rank < cards[i-1].Rank {
			return false
		}
		if cards[i].Rank == cards[i-1].Rank {
			if cards[i].Suit < cards[i-1].Suit {
				return false
			}
		}
	}
	return true
}
