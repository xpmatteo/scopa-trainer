package application

import (
	"sort"

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

	// Get player hand and sort it by rank first, then by suit
	playerHand := s.gameState.Deck.CardsAt(domain.AIHandLocation)
	sort.Slice(playerHand, func(i, j int) bool {
		// First compare by rank
		if playerHand[i].Rank != playerHand[j].Rank {
			return playerHand[i].Rank < playerHand[j].Rank
		}
		// If ranks are equal, compare by suit
		return string(playerHand[i].Suit) < string(playerHand[j].Suit)
	})

	// Update the UI model
	s.model.ShowNewGameButton = false
	s.model.TableCards = s.gameState.Deck.CardsAt(domain.TableLocation)
	s.model.PlayerHand = playerHand
	s.model.GameInProgress = true
	s.model.PlayerTurn = s.gameState.PlayerTurn

	if s.model.PlayerTurn {
		s.model.GamePrompt = "Your turn. Select a card to play."
	} else {
		s.model.GamePrompt = "AI is thinking..."
	}

	return s.model
}
