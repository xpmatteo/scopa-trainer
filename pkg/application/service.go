package application

import (
	"sort"

	"github.com/xpmatteo/scopa-trainer/pkg/domain"
)

// GameService handles the application logic for the game
type GameService struct {
	gameState    *domain.GameState
	selectedCard domain.Card
}

// NewGameService creates a new game service with initial state
func NewGameService() *GameService {
	return &GameService{
		gameState:    nil,
		selectedCard: domain.NO_CARD_SELECTED,
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
	model.PlayerHand = sortCards(s.gameState.Deck.CardsAt(domain.PlayerHandLocation))
	model.GameInProgress = true
	model.PlayerTurn = s.gameState.PlayerTurn
	model.SelectedCard = s.selectedCard

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
	s.selectedCard = s.gameState.Deck.CardsAt(domain.PlayerHandLocation)[0]

	// Return the UI model
	return s.GetUIModel()
}

// sortCards sorts the cards by rank and suit
func sortCards(cards []domain.Card) []domain.Card {
	sort.Slice(cards, func(i, j int) bool {
		// First compare by rank
		if cards[i].Rank != cards[j].Rank {
			return cards[i].Rank < cards[j].Rank
		}
		// If ranks are equal, compare by suit
		return string(cards[i].Suit) < string(cards[j].Suit)
	})
	return cards
}

// SelectCard handles the selection of a card from the player's hand
func (s *GameService) SelectCard(suit domain.Suit, rank domain.Rank) domain.UIModel {
	// Find the card in the player's hand
	for _, card := range s.gameState.Deck.CardsAt(domain.PlayerHandLocation) {
		if card.Suit == suit && card.Rank == rank {
			// If the card is already selected, deselect it
			if s.selectedCard == card {
				s.selectedCard = domain.NO_CARD_SELECTED
			} else {
				// Otherwise, select it
				s.selectedCard = card
			}
			break
		}
	}

	return s.GetUIModel()
}
