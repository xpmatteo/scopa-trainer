package application

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/xpmatteo/scopa-trainer/pkg/domain"
)

// GameService handles the application logic for the game
type GameService struct {
	gameState          *domain.GameState
	selectedCard       domain.Card
	selectedTableCards []domain.Card // Cards selected for combination capture
}

// NewGameService creates a new game service with initial state
func NewGameService() *GameService {
	return &GameService{
		gameState:          nil,
		selectedCard:       domain.NO_CARD_SELECTED,
		selectedTableCards: []domain.Card{},
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
	model.SelectedTableCards = s.selectedTableCards

	// Set card counts and capture cards
	model.DeckCount = len(s.gameState.Deck.CardsAt(domain.DeckLocation))
	model.PlayerCaptureCards = sortCards(s.gameState.Deck.CardsAt(domain.PlayerCapturesLocation))
	model.PlayerCaptureCount = len(model.PlayerCaptureCards)
	model.AICaptureCards = sortCards(s.gameState.Deck.CardsAt(domain.AICapturesLocation))
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
			// Get all possible capture options
			captureOptions := s.findCaptureOptions(s.selectedCard)
			model.CaptureOptions = captureOptions

			// Check if we can confirm the current capture selection
			model.CanConfirmCapture = s.canCompleteCaptureSelection()

			if len(captureOptions) > 0 {
				// We have capture options
				if len(captureOptions) == 1 && len(captureOptions[0]) == 1 {
					// Single card capture
					model.GamePrompt = "Click on the matching card to capture it, or select a different card."
				} else if len(s.selectedTableCards) > 0 {
					// Some table cards are selected for combination capture
					if model.CanConfirmCapture {
						model.GamePrompt = "Click 'Confirm Capture' to complete the capture, or select different cards."
					} else {
						sum := 0
						for _, card := range s.selectedTableCards {
							sum += card.Value()
						}
						model.GamePrompt = fmt.Sprintf("Selected cards sum to %d (need %d). Select more cards or click a selected card to deselect it.", sum, s.selectedCard.Value())
					}
				} else {
					// Combination capture, nothing selected yet
					model.GamePrompt = "Select table cards that sum to " + strconv.Itoa(s.selectedCard.Value()) + "."
				}
				// Set CanPlaySelectedCard to false when capture is possible
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

	// Get all possible capture options
	options := s.findCaptureOptions(card)

	// If there are any options, a capture is possible
	return len(options) > 0
}

// isValidCaptureSelection checks if the currently selected table cards form a valid capture
// with the selected hand card
func (s *GameService) isValidCaptureSelection() bool {
	// If no card is selected from hand, or no table cards selected, not valid
	if s.selectedCard == domain.NO_CARD_SELECTED || len(s.selectedTableCards) == 0 {
		return false
	}

	// If just one card is selected, it must match the rank of the hand card
	if len(s.selectedTableCards) == 1 && s.selectedTableCards[0].Rank == s.selectedCard.Rank {
		return true
	}

	// For multiple cards, calculate the sum of their values
	sum := 0
	for _, card := range s.selectedTableCards {
		sum += card.Value()
	}

	// The sum must equal the hand card's value
	return sum == s.selectedCard.Value()
}

// toggleTableCardSelection adds or removes a card from the selectedTableCards slice
func (s *GameService) toggleTableCardSelection(card domain.Card) {
	// Check if the card is already selected
	for i, selectedCard := range s.selectedTableCards {
		if selectedCard == card {
			// Card is already selected, so remove it
			s.selectedTableCards = append(s.selectedTableCards[:i], s.selectedTableCards[i+1:]...)
			return
		}
	}

	// If we've gotten here, card was not already selected, so add it
	s.selectedTableCards = append(s.selectedTableCards, card)
}

// canCompleteCaptureSelection checks if we have a valid capture selection that can be confirmed
func (s *GameService) canCompleteCaptureSelection() bool {
	return s.isValidCaptureSelection()
}

// findCaptureOptions returns all possible card combinations that can be captured
// Priority order:
// 1. Single card with the same rank
// 2. Multiple cards whose values sum to the played card's value
func (s *GameService) findCaptureOptions(card domain.Card) [][]domain.Card {
	if card == domain.NO_CARD_SELECTED {
		return nil
	}

	tableCards := s.gameState.Deck.CardsAt(domain.TableLocation)

	// First, check for single card matches (these take precedence)
	for _, tableCard := range tableCards {
		if tableCard.Rank == card.Rank {
			// Return only this single card match
			return [][]domain.Card{{tableCard}}
		}
	}

	// If no single card match, find all combinations of table cards that sum to card value
	cardValue := card.Value()
	return s.findAllCombinations(tableCards, cardValue)
}

// findAllCombinations returns all combinations of cards that sum to the target value
func (s *GameService) findAllCombinations(cards []domain.Card, target int) [][]domain.Card {
	var result [][]domain.Card

	// Try combinations of different sizes (2 to N cards)
	for size := 2; size <= len(cards); size++ {
		combinations := s.generateCombinations(cards, size)
		for _, combo := range combinations {
			sum := 0
			for _, c := range combo {
				sum += c.Value()
			}

			// If the sum matches the target, add this combination to results
			if sum == target {
				result = append(result, combo)
			}
		}
	}

	return result
}

// generateCombinations returns all possible combinations of k cards from the input slice
func (s *GameService) generateCombinations(cards []domain.Card, k int) [][]domain.Card {
	var result [][]domain.Card
	n := len(cards)

	// Base cases
	if k > n {
		return result
	}

	if k == 1 {
		// Each card is a combination of size 1
		for _, card := range cards {
			result = append(result, []domain.Card{card})
		}
		return result
	}

	// Generate combinations recursively
	for i := 0; i <= n-k; i++ {
		// Take the current card
		current := cards[i]

		// Generate combinations for remaining cards
		subCombinations := s.generateCombinations(cards[i+1:], k-1)

		// Add current card to each sub-combination
		for _, subCombo := range subCombinations {
			combo := append([]domain.Card{current}, subCombo...)
			result = append(result, combo)
		}
	}

	return result
}

// CaptureCombination captures a combination of cards from the table
func (s *GameService) CaptureCombination(tableCards []domain.Card) {
	if s.selectedCard == domain.NO_CARD_SELECTED || len(tableCards) == 0 {
		return
	}

	// Verify all cards are on the table
	for _, card := range tableCards {
		if s.gameState.Deck.GetCardLocation(card) != domain.TableLocation {
			return
		}
	}

	// Move the selected card from hand to capture pile
	s.gameState.Deck.MoveCard(s.selectedCard, domain.PlayerHandLocation, domain.PlayerCapturesLocation)

	// Move all table cards in the combination to the capture pile
	for _, card := range tableCards {
		s.gameState.Deck.MoveCard(card, domain.TableLocation, domain.PlayerCapturesLocation)
	}

	// Clear the selected card and table card selection
	s.selectedCard = domain.NO_CARD_SELECTED
	s.selectedTableCards = []domain.Card{}

	// Switch turn to AI
	s.gameState.Status = domain.StatusAITurn

	// Check if new cards need to be dealt
	s.DealNewCardsIfNeeded()
}

// ConfirmCapture confirms the current capture selection
func (s *GameService) ConfirmCapture() {
	// Verify that we have a valid capture selection
	if !s.isValidCaptureSelection() {
		return
	}

	// Use the existing CaptureCombination method with our selected table cards
	s.CaptureCombination(s.selectedTableCards)
}

// StartNewGame initializes a new game and returns the updated UI model
func (s *GameService) StartNewGame() {
	// Initialize a new game state
	gameState := domain.NewGameState()
	s.gameState = &gameState
	s.selectedCard = domain.NO_CARD_SELECTED
	s.selectedTableCards = []domain.Card{}
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
		// Check if the ranks match for capture (direct match - priority rule)
		if clickedCard.Rank == s.selectedCard.Rank {
			// Direct rank match takes precedence - capture immediately
			s.CaptureCombination([]domain.Card{clickedCard})
			return
		}

		// For combination captures, toggle the selection of this card
		options := s.findCaptureOptions(s.selectedCard)
		if len(options) > 0 {
			// Toggle this card's selection
			s.toggleTableCardSelection(clickedCard)

			// If the selection is now valid, user can confirm it
			// The UI will show a confirmation button
			return
		}

		// If no valid capture options, keep the hand card selected
		return
	}

	// If table card clicked without a hand card selected, do nothing
	if isTableCard && s.selectedCard == domain.NO_CARD_SELECTED {
		// Don't change the selected card
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

	// Find all possible capture options for this card
	options := s.findAICaptureOptions(aiCard)

	if len(options) > 0 {
		// AI has at least one capture option
		// For simplicity, always choose the first option
		captureCards := options[0]

		// Move AI card to captures
		s.gameState.Deck.MoveCard(aiCard, domain.AIHandLocation, domain.AICapturesLocation)

		// Move all table cards in the combination to AI captures
		for _, card := range captureCards {
			s.gameState.Deck.MoveCard(card, domain.TableLocation, domain.AICapturesLocation)
		}
	} else {
		// No captures possible, play card to table
		s.gameState.Deck.MoveCard(aiCard, domain.AIHandLocation, domain.TableLocation)
	}

	// Switch turn to player
	s.gameState.Status = domain.StatusPlayerTurn

	// Check if new cards need to be dealt
	s.DealNewCardsIfNeeded()
}

// findAICaptureOptions works like findCaptureOptions but for the AI player
func (s *GameService) findAICaptureOptions(card domain.Card) [][]domain.Card {
	if card == domain.NO_CARD_SELECTED {
		return nil
	}

	tableCards := s.gameState.Deck.CardsAt(domain.TableLocation)

	// First, check for single card matches (these take precedence)
	for _, tableCard := range tableCards {
		if tableCard.Rank == card.Rank {
			// Return only this single card match
			return [][]domain.Card{{tableCard}}
		}
	}

	// If no single card match, find all combinations of table cards that sum to card value
	cardValue := card.Value()
	return s.findAllCombinations(tableCards, cardValue)
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

// SetupCombinationTest sets up a test scenario for combination captures
// It creates a game state with specific cards for testing the capture combinations:
// - Cards on the table: Ranks 1, 2, 3, 4 of different suits
// - Card in player's hand: Rank 7
// This allows testing of capturing combinations like 3+4=7 or 1+2+4=7
func (s *GameService) SetupCombinationTest() {
	// Start with a new game to initialize the game state
	s.StartNewGame()

	// Clear the table
	tableCards := s.gameState.Deck.CardsAt(domain.TableLocation)
	for _, card := range tableCards {
		s.gameState.Deck.MoveCard(card, domain.TableLocation, domain.DeckLocation)
	}

	// Clear the player's hand
	playerCards := s.gameState.Deck.CardsAt(domain.PlayerHandLocation)
	for _, card := range playerCards {
		s.gameState.Deck.MoveCard(card, domain.PlayerHandLocation, domain.DeckLocation)
	}

	// Clear the AI's hand
	aiCards := s.gameState.Deck.CardsAt(domain.AIHandLocation)
	for _, card := range aiCards {
		s.gameState.Deck.MoveCard(card, domain.AIHandLocation, domain.DeckLocation)
	}

	// Find and move cards to the table
	s.gameState.Deck.MoveCardMatching(domain.DeckLocation, domain.TableLocation, domain.Asso, domain.Coppe)
	s.gameState.Deck.MoveCardMatching(domain.DeckLocation, domain.TableLocation, domain.Due, domain.Bastoni)
	s.gameState.Deck.MoveCardMatching(domain.DeckLocation, domain.TableLocation, domain.Tre, domain.Spade)
	s.gameState.Deck.MoveCardMatching(domain.DeckLocation, domain.TableLocation, domain.Quattro, domain.Denari)

	// Add a card with rank 7 to player's hand
	s.gameState.Deck.MoveCardMatching(domain.DeckLocation, domain.PlayerHandLocation, domain.Sette, domain.Denari)

	// Set game state to player's turn
	s.gameState.Status = domain.StatusPlayerTurn
}
