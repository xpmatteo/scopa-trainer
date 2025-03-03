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

func (s *NewGameTestSuite) TestInitialUIModelHasNewGameButton() {
	// When the application starts
	model := s.service.GetUIModel()

	// Then the UI model should have a "New Game" button
	assert.True(s.T(), model.ShowNewGameButton, "New Game button should be displayed initially")
}

func (s *NewGameTestSuite) TestStartNewGame() {
	// Given the application has started
	initialModel := s.service.GetUIModel()
	assert.True(s.T(), initialModel.ShowNewGameButton, "New Game button should be displayed initially")

	// When the player clicks the "New Game" button
	updatedModel := s.service.StartNewGame()

	assert.Equal(s.T(), 0, len(updatedModel.TableCards), "Table should have no cards")
	assert.Equal(s.T(), 10, len(updatedModel.PlayerHand), "Player should have 10 cards")
	assert.False(s.T(), updatedModel.ShowNewGameButton, "New Game button should be hidden after starting a game")
}
