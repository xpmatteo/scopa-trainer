package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
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
	updatedModel := service.StartNewGame()

	assert.Equal(t, 0, len(updatedModel.TableCards), "Table should have no cards")
	assert.Equal(t, 10, len(updatedModel.PlayerHand), "Player should have 10 cards")
	assert.False(t, updatedModel.ShowNewGameButton, "New Game button should be hidden after starting a game")
}
