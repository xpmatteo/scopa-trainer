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
	model.GameInProgress = s.gameState.Status != domain.StatusGameNotStarted
	model.PlayerTurn = s.gameState.Status == domain.StatusPlayerTurn
	model.SelectedCard = s.selectedCard

	// Set card counts and capture cards
	model.DeckCount = len(s.gameState.Deck.CardsAt(domain.DeckLocation))
	model.PlayerCaptureCards = s.gameState.Deck.CardsAt(domain.PlayerCapturesLocation)
	model.PlayerCaptureCount = len(model.PlayerCaptureCards)
	model.AICaptureCards = s.gameState.Deck.CardsAt(domain.AICapturesLocation)
	model.AICaptureCount = len(model.AICaptureCards)

	// Calculate the score (updated continuously)
	model.Score = domain.CalculateScore(model.PlayerCaptureCards, model.AICaptureCards)

	// Check if the game is over
	model.GameOver = s.gameState.Status == domain.StatusGameOver
	if model.GameOver {
		model.ShowNewGameButton = true
		model.GamePrompt = "Game Over! Check out your score and the AI's score."
	} else if model.PlayerTurn {
		// Player's turn
		if s.selectedCard == domain.NO_CARD_SELECTED {
			model.GamePrompt = "Select a card from your hand to play."
			model.CanPlaySelectedCard = false
		} else {
			// Check if the selected card can capture any cards
			if s.canCaptureAnyCard(s.selectedCard) {
				model.GamePrompt = "Click on the table card(s) you want to capture, or select a different card."
				// Set CanPlaySelectedCard to false when capture is possible to maintain compatibility with existing tests
				model.CanPlaySelectedCard = false
			} else {
				model.GamePrompt = "This card cannot capture any cards. Click on the table to discard it, or select a different card."
				model.CanPlaySelectedCard = true
			}
		}
	} else {
		// AI's turn
		model.GamePrompt = "AI is thinking..."
	}

	return model
}

// canCaptureAnyCard checks if the given card can capture any card on the table
func (s *GameService) canCaptureAnyCard(card domain.Card) bool {
	// If no card is selected, no capture is possible
	if card == domain.NO_CARD_SELECTED {
		return false
	}

	// Check if any table card has the same rank as the selected card
	for _, tableCard := range s.gameState.Deck.CardsAt(domain.TableLocation) {
		if tableCard.Rank == card.Rank {
			return true
		}
	}

	return false
}

// StartNewGame initializes a new game and returns the updated UI model
func (s *GameService) StartNewGame() {
	// Initialize a new game state
	gameState := domain.NewGameState()
	s.gameState = &gameState
	s.selectedCard = domain.NO_CARD_SELECTED
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
func (s *GameService) SelectCard(suit domain.Suit, rank domain.Rank) {
	// Find the card that was clicked
	clickedCard := domain.Card{
		Suit: suit,
		Rank: rank,
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

			// Switch turn to AI
			s.gameState.Status = domain.StatusAITurn

			// Check if new cards need to be dealt
			s.DealNewCardsIfNeeded()
		}
		// If ranks don't match, keep the hand card selected
		return
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
}

// PlaySelectedCard moves the currently selected card from the player's hand to the table
func (s *GameService) PlaySelectedCard() {
	// Check if a card is selected
	if s.selectedCard == domain.NO_CARD_SELECTED {
		return
	}

	// Check if a capture is possible
	if s.canCaptureAnyCard(s.selectedCard) {
		// Cannot play to table if capture is possible
		return
	}

	// Move the selected card to the table
	s.gameState.Deck.MoveCard(s.selectedCard, domain.PlayerHandLocation, domain.TableLocation)

	// Clear the selected card
	s.selectedCard = domain.NO_CARD_SELECTED

	// Switch turn to AI
	s.gameState.Status = domain.StatusAITurn

	// Check if new cards need to be dealt
	s.DealNewCardsIfNeeded()
}

// DealNewCardsIfNeeded checks if hands are empty and deals new cards if needed
// Returns true if new cards were dealt
func (s *GameService) DealNewCardsIfNeeded() bool {
	// If no game in progress, do nothing
	if s.gameState == nil {
		return false
	}

	playerHand := s.gameState.Deck.CardsAt(domain.PlayerHandLocation)
	aiHand := s.gameState.Deck.CardsAt(domain.AIHandLocation)
	deckCards := s.gameState.Deck.CardsAt(domain.DeckLocation)

	// Check if both hands are empty and there are cards in the deck
	if len(playerHand) == 0 && len(aiHand) == 0 {
		if len(deckCards) > 0 {
			// Calculate how many cards to deal to each player
			cardsPerPlayer := 10
			if len(deckCards) < 20 {
				// If fewer than 20 cards, distribute evenly
				cardsPerPlayer = len(deckCards) / 2
			}

			// Deal cards to each player
			s.gameState.Deck.DealCards(domain.DeckLocation, domain.PlayerHandLocation, cardsPerPlayer)
			s.gameState.Deck.DealCards(domain.DeckLocation, domain.AIHandLocation, cardsPerPlayer)

			return true
		} else {
			// Both hands are empty and deck is empty, game is over
			s.gameState.Status = domain.StatusGameOver
		}
	}

	return false
}

// PlayAITurn handles the AI's turn
func (s *GameService) PlayAITurn() {
	// Check if it's the AI's turn
	if s.gameState.Status != domain.StatusAITurn {
		return
	}

	// Get the first card from AI's hand
	aiCards := s.gameState.Deck.CardsAt(domain.AIHandLocation)
	if len(aiCards) == 0 {
		// No cards in AI hand, nothing to do
		return
	}

	// Select the first card
	aiCard := aiCards[0]

	// Check if the card can capture any card on the table
	tableCards := s.gameState.Deck.CardsAt(domain.TableLocation)
	captured := false

	for _, tableCard := range tableCards {
		if tableCard.Rank == aiCard.Rank {
			// Capture the card
			s.gameState.Deck.MoveCard(aiCard, domain.AIHandLocation, domain.AICapturesLocation)
			s.gameState.Deck.MoveCard(tableCard, domain.TableLocation, domain.AICapturesLocation)
			captured = true
			break // Only capture the first matching card
		}
	}

	// If no capture was made, play the card to the table
	if !captured {
		s.gameState.Deck.MoveCard(aiCard, domain.AIHandLocation, domain.TableLocation)
	}

	// Switch turn to player
	s.gameState.Status = domain.StatusPlayerTurn

	// Check if new cards need to be dealt
	s.DealNewCardsIfNeeded()
}

// SetGameOver sets the game state to game over
func (s *GameService) SetGameOver() {
	if s.gameState != nil {
		s.gameState.Status = domain.StatusGameOver
	}
}

// GetDeck returns the current game deck
// Note: This method is primarily for testing purposes
func (s *GameService) GetDeck() *domain.Deck {
	if s.gameState == nil {
		return nil
	}
	return s.gameState.Deck
}
