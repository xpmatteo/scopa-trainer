package application

import (
	"sort"

	"github.com/xpmatteo/scopa-trainer/pkg/domain"
)

// GameService handles the application logic for the game
type GameService struct {
	gameState *domain.GameState
}

// NewGameService creates a new game service with initial state
func NewGameService() *GameService {
	return &GameService{
		gameState: nil,
	}
}

// GetUIModel returns the current UI model
func (s *GameService) GetUIModel() domain.UIModel {
	if s.gameState == nil {
		// Return initial UI model when no game is in progress
		return domain.NewUIModel()
	}

	// Generate UI model based on current game state
	model := domain.NewUIModel()
	model.ShowNewGameButton = false
	model.TableCards = s.gameState.Deck.CardsAt(domain.TableLocation)
	model.PlayerHand = s.getSortedPlayerHand()
	model.GameInProgress = true
	model.PlayerTurn = s.gameState.PlayerTurn

	if model.PlayerTurn {
		model.GamePrompt = "Your turn. Select a card to play."
	} else {
		model.GamePrompt = "AI is thinking..."
	}

	return model
}

// StartNewGame initializes a new game and returns the updated UI model
func (s *GameService) StartNewGame() domain.UIModel {
	// Initialize a new game state
	gameState := domain.NewGameState()
	s.gameState = &gameState

	// Return the UI model
	return s.GetUIModel()
}

// getSortedPlayerHand returns the player's hand sorted by rank and suit
func (s *GameService) getSortedPlayerHand() []domain.Card {
	playerHand := s.gameState.Deck.CardsAt(domain.AIHandLocation)
	sort.Slice(playerHand, func(i, j int) bool {
		// First compare by rank
		if playerHand[i].Rank != playerHand[j].Rank {
			return playerHand[i].Rank < playerHand[j].Rank
		}
		// If ranks are equal, compare by suit
		return string(playerHand[i].Suit) < string(playerHand[j].Suit)
	})
	return playerHand
}
