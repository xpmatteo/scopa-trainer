package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xpmatteo/scopa-trainer/pkg/domain"
)

func TestFindCaptureOptions(t *testing.T) {
	// Arrange
	service := NewGameService()
	service.StartNewGame()
	
	// Clear the table
	clearTable(service)
	
	// Place specific cards on the table
	addCardToTable(service, domain.Card{Suit: domain.Coppe, Rank: domain.Asso})    // 1
	addCardToTable(service, domain.Card{Suit: domain.Bastoni, Rank: domain.Due})   // 2
	addCardToTable(service, domain.Card{Suit: domain.Spade, Rank: domain.Tre})     // 3
	addCardToTable(service, domain.Card{Suit: domain.Denari, Rank: domain.Quattro}) // 4
	
	// Add a rank 7 card to player's hand
	handCard := domain.Card{Suit: domain.Coppe, Rank: domain.Sette} // 7
	addCardToHand(service, handCard)
	
	// Act
	options := service.findCaptureOptions(handCard)
	
	// Assert
	assert.NotEmpty(t, options)
	
	// Check for 3+4=7 combination
	hasThreePlusFour := false
	for _, option := range options {
		if len(option) == 2 {
			hasRanks := hasCardWithRank(option, domain.Tre) && hasCardWithRank(option, domain.Quattro)
			if hasRanks {
				hasThreePlusFour = true
				break
			}
		}
	}
	assert.True(t, hasThreePlusFour, "Should find the 3+4=7 combination")
	
	// Check for 1+2+4=7 combination
	hasOnePlusTwoPlusFour := false
	for _, option := range options {
		if len(option) == 3 {
			hasRanks := hasCardWithRank(option, domain.Asso) && 
				hasCardWithRank(option, domain.Due) && 
				hasCardWithRank(option, domain.Quattro)
			if hasRanks {
				hasOnePlusTwoPlusFour = true
				break
			}
		}
	}
	assert.True(t, hasOnePlusTwoPlusFour, "Should find the 1+2+4=7 combination")
}

func TestSingleCardCaptureTakesPrecedence(t *testing.T) {
	// Arrange
	service := NewGameService()
	service.StartNewGame()
	
	// Clear the table
	clearTable(service)
	
	// Place specific cards on the table
	addCardToTable(service, domain.Card{Suit: domain.Coppe, Rank: domain.Cinque})  // 5
	addCardToTable(service, domain.Card{Suit: domain.Bastoni, Rank: domain.Due})   // 2
	addCardToTable(service, domain.Card{Suit: domain.Spade, Rank: domain.Tre})     // 3
	
	// Add a rank 5 card to player's hand
	handCard := domain.Card{Suit: domain.Bastoni, Rank: domain.Cinque} // 5
	addCardToHand(service, handCard)
	
	// Act
	options := service.findCaptureOptions(handCard)
	
	// Assert
	assert.Equal(t, 1, len(options), "Should only find one option when a single card match exists")
	assert.Equal(t, 1, len(options[0]), "Should match only the single card with same rank")
	assert.Equal(t, domain.Cinque, options[0][0].Rank, "Should match the card with rank 5")
}

func TestCaptureCombination(t *testing.T) {
	// Arrange
	service := NewGameService()
	service.StartNewGame()
	
	// Clear the table
	clearTable(service)
	
	// Place specific cards on the table
	tableCard1 := domain.Card{Suit: domain.Bastoni, Rank: domain.Due}   // 2
	tableCard2 := domain.Card{Suit: domain.Spade, Rank: domain.Tre}     // 3
	addCardToTable(service, tableCard1)
	addCardToTable(service, tableCard2)
	
	// Add a rank 5 card to player's hand and select it
	handCard := domain.Card{Suit: domain.Bastoni, Rank: domain.Cinque}  // 5
	addCardToHand(service, handCard)
	service.selectedCard = handCard
	
	// Verify initial state
	assert.Equal(t, 2, len(service.gameState.Deck.CardsAt(domain.TableLocation)))
	assert.Equal(t, 0, len(service.gameState.Deck.CardsAt(domain.PlayerCapturesLocation)))
	
	// Act
	service.CaptureCombination([]domain.Card{tableCard1, tableCard2})
	
	// Assert
	assert.Equal(t, 0, len(service.gameState.Deck.CardsAt(domain.TableLocation)), 
		"Table should be empty after capture")
	assert.Equal(t, 3, len(service.gameState.Deck.CardsAt(domain.PlayerCapturesLocation)), 
		"Player should have 3 cards in capture pile")
	assert.Equal(t, domain.NO_CARD_SELECTED, service.selectedCard, 
		"No card should be selected after capture")
	assert.Equal(t, domain.StatusAITurn, service.gameState.Status, 
		"Should be AI's turn after capture")
}

func TestCannotPlayCardWhenCombinationCaptureIsPossible(t *testing.T) {
	// Arrange
	service := NewGameService()
	service.StartNewGame()
	
	// Clear the table
	clearTable(service)
	
	// Place specific cards on the table
	addCardToTable(service, domain.Card{Suit: domain.Bastoni, Rank: domain.Due})   // 2
	addCardToTable(service, domain.Card{Suit: domain.Spade, Rank: domain.Tre})     // 3
	initialTableCount := len(service.gameState.Deck.CardsAt(domain.TableLocation))
	
	// Add a rank 5 card to player's hand and select it
	handCard := domain.Card{Suit: domain.Bastoni, Rank: domain.Cinque}  // 5
	addCardToHand(service, handCard)
	service.selectedCard = handCard
	initialHandCount := len(service.gameState.Deck.CardsAt(domain.PlayerHandLocation))
	
	// Get UI model first to verify capture is not allowed
	model := service.GetUIModel()
	assert.False(t, model.CanPlaySelectedCard, "Should not allow playing card when capture is possible")
	
	// Act
	service.PlaySelectedCard()
	
	// Assert
	assert.Equal(t, initialTableCount, len(service.gameState.Deck.CardsAt(domain.TableLocation)), 
		"Table should remain unchanged")
	assert.Equal(t, initialHandCount, len(service.gameState.Deck.CardsAt(domain.PlayerHandLocation)), 
		"Hand should remain unchanged")
	assert.Equal(t, domain.StatusPlayerTurn, service.gameState.Status, 
		"Should still be player's turn")
	assert.Equal(t, handCard, service.selectedCard, 
		"Card should remain selected")
}

func TestAIPlayerWithCombinationCapture(t *testing.T) {
	// Arrange
	service := NewGameService()
	service.StartNewGame()
	
	// Clear the table and AI hand
	clearTable(service)
	clearAIHand(service)
	
	// Place specific cards on the table
	addCardToTable(service, domain.Card{Suit: domain.Bastoni, Rank: domain.Due})   // 2
	addCardToTable(service, domain.Card{Suit: domain.Spade, Rank: domain.Tre})     // 3
	
	// Add a rank 5 card to AI's hand
	aiCard := domain.Card{Suit: domain.Bastoni, Rank: domain.Cinque}  // 5
	addCardToAIHand(service, aiCard)
	
	// Set game status to AI turn
	service.gameState.Status = domain.StatusAITurn
	
	// Act
	service.PlayAITurn()
	
	// Assert
	assert.Equal(t, 0, len(service.gameState.Deck.CardsAt(domain.TableLocation)), 
		"Table should be empty after capture")
	assert.Equal(t, 0, len(service.gameState.Deck.CardsAt(domain.AIHandLocation)), 
		"AI hand should be empty after playing")
	assert.Equal(t, 3, len(service.gameState.Deck.CardsAt(domain.AICapturesLocation)), 
		"AI should have 3 cards in capture pile")
	assert.Equal(t, domain.StatusPlayerTurn, service.gameState.Status, 
		"Should be player's turn after AI capture")
}

// Helper functions for test setup

func clearTable(service *GameService) {
	tableCards := service.gameState.Deck.CardsAt(domain.TableLocation)
	for _, card := range tableCards {
		service.gameState.Deck.MoveCard(card, domain.TableLocation, domain.DeckLocation)
	}
}

func clearAIHand(service *GameService) {
	aiCards := service.gameState.Deck.CardsAt(domain.AIHandLocation)
	for _, card := range aiCards {
		service.gameState.Deck.MoveCard(card, domain.AIHandLocation, domain.DeckLocation)
	}
}

func addCardToTable(service *GameService, card domain.Card) {
	// Find the card in the deck and move it to the table
	deckCards := service.gameState.Deck.CardsAt(domain.DeckLocation)
	for _, deckCard := range deckCards {
		if deckCard.Rank == card.Rank && deckCard.Suit == card.Suit {
			service.gameState.Deck.MoveCard(deckCard, domain.DeckLocation, domain.TableLocation)
			return
		}
	}
	
	// If card not found, create a custom deck with the card
	setupCustomDeck(service, card, domain.TableLocation)
}

func addCardToHand(service *GameService, card domain.Card) {
	// Find the card in the deck and move it to the hand
	deckCards := service.gameState.Deck.CardsAt(domain.DeckLocation)
	for _, deckCard := range deckCards {
		if deckCard.Rank == card.Rank && deckCard.Suit == card.Suit {
			service.gameState.Deck.MoveCard(deckCard, domain.DeckLocation, domain.PlayerHandLocation)
			return
		}
	}
	
	// If card not found, create a custom deck with the card
	setupCustomDeck(service, card, domain.PlayerHandLocation)
}

func addCardToAIHand(service *GameService, card domain.Card) {
	// Find the card in the deck and move it to the AI hand
	deckCards := service.gameState.Deck.CardsAt(domain.DeckLocation)
	for _, deckCard := range deckCards {
		if deckCard.Rank == card.Rank && deckCard.Suit == card.Suit {
			service.gameState.Deck.MoveCard(deckCard, domain.DeckLocation, domain.AIHandLocation)
			return
		}
	}
	
	// If card not found, create a custom deck with the card
	setupCustomDeck(service, card, domain.AIHandLocation)
}

// Helper to ensure a specific card exists at a specific location
func setupCustomDeck(service *GameService, card domain.Card, location domain.Location) {
	// For testing, check all cards in the deck
	// This is a workaround since we can't directly access unexported fields
	for _, suit := range domain.AllSuits() {
		for _, rank := range domain.AllRanks() {
			testCard := domain.Card{Suit: suit, Rank: rank}
			
			// If this is our target card, ensure it's at the right location
			if testCard.Rank == card.Rank && testCard.Suit == card.Suit {
				// Move it from wherever it is to the desired location
				currentLocation := service.gameState.Deck.GetCardLocation(testCard)
				service.gameState.Deck.MoveCard(testCard, currentLocation, location)
				return
			}
		}
	}
}

func hasCardWithRank(cards []domain.Card, rank domain.Rank) bool {
	for _, card := range cards {
		if card.Rank == rank {
			return true
		}
	}
	return false
}