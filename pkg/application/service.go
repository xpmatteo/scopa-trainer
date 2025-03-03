package application

import (
	"github.com/xpmatteo/scopa-trainer/pkg/domain"
)

// GameService handles the application logic for the game
type GameService struct {
	model     domain.UIModel
	gameState *domain.GameState
}

// NewGameService creates a new game service with initial state
func NewGameService() *GameService {
	return &GameService{
		model:     domain.NewUIModel(),
		gameState: nil,
	}
}

// GetUIModel returns the current UI model
func (s *GameService) GetUIModel() domain.UIModel {
	return s.model
}

// StartNewGame initializes a new game and returns the updated UI model
func (s *GameService) StartNewGame() domain.UIModel {
	// Initialize a new game state
	gameState := domain.NewGameState()
	s.gameState = &gameState

	// Update the UI model
	s.model.ShowNewGameButton = false
	s.model.TableCards = s.gameState.TableCards
	s.model.PlayerHand = s.gameState.PlayerHand
	s.model.GameInProgress = true
	s.model.PlayerTurn = s.gameState.PlayerTurn

	if s.model.PlayerTurn {
		s.model.GamePrompt = "Your turn. Select a card to play."
	} else {
		s.model.GamePrompt = "AI is thinking..."
	}

	return s.model
}
