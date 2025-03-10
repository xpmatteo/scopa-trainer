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
	s.selectedCard = domain.NO_CARD_SELECTED

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

// SelectCard handles the selection of a card from the player's hand or capturing a card from the table
func (s *GameService) SelectCard(suit domain.Suit, rank domain.Rank) domain.UIModel {
	// Find the card that was clicked
	clickedCard := domain.Card{
		Suit:  suit,
		Rank:  rank,
		Name:  rank.String(),
		Value: rank.Value(),
	}

	// Check if the clicked card is on the table
	isTableCard := false
	for _, card := range s.gameState.Deck.CardsAt(domain.TableLocation) {
		if card.Suit == suit && card.Rank == rank {
			isTableCard = true
			clickedCard = card // Use the actual card from the table
			break
		}
	}

	// Check if the clicked card is in the player's hand
	isHandCard := false
	for _, card := range s.gameState.Deck.CardsAt(domain.PlayerHandLocation) {
		if card.Suit == suit && card.Rank == rank {
			isHandCard = true
			clickedCard = card // Use the actual card from the hand
			break
		}
	}

	// If a card from the table was clicked and we have a selected hand card
	if isTableCard && s.selectedCard != domain.NO_CARD_SELECTED {
		// Check if the ranks match for capture
		if clickedCard.Rank == s.selectedCard.Rank {
			// Capture the card
			s.gameState.Deck.MoveCard(s.selectedCard, domain.PlayerHandLocation, domain.PlayerCapturesLocation)
			s.gameState.Deck.MoveCard(clickedCard, domain.TableLocation, domain.PlayerCapturesLocation)

			// Clear the selected card
			s.selectedCard = domain.NO_CARD_SELECTED
		}
		// If ranks don't match, keep the hand card selected
		return s.GetUIModel()
	}

	// If a card from the hand was clicked
	if isHandCard {
		// If the card is already selected, deselect it
		if s.selectedCard == clickedCard {
			s.selectedCard = domain.NO_CARD_SELECTED
		} else {
			// Otherwise, select it
			s.selectedCard = clickedCard
		}
	}
	// If a table card was clicked without a hand card selected, do nothing
	// The selected card remains unchanged

	return s.GetUIModel()
}
