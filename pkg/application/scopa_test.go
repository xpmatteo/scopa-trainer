package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xpmatteo/scopa-trainer/pkg/domain"
)

func TestPlayerScopa(t *testing.T) {
	// Setup
	service := NewGameService()
	service.StartNewGame()

	// Set up for a clean test
	service.SetupCombinationTest() // This creates a fresh game with specific cards

	// Keep only one card on the table and one in the player's hand to make testing simpler
	// First, clear the table
	tableCards := service.gameState.Deck.CardsAt(domain.TableLocation)
	for _, card := range tableCards {
		service.gameState.Deck.MoveCard(card, domain.TableLocation, domain.DeckLocation)
	}

	// Add one card to the table (Asso di Coppe)
	service.gameState.Deck.MoveCardMatching(domain.DeckLocation, domain.TableLocation, domain.Asso, domain.Coppe)

	// Clear player hand
	playerCards := service.gameState.Deck.CardsAt(domain.PlayerHandLocation)
	for _, card := range playerCards {
		service.gameState.Deck.MoveCard(card, domain.PlayerHandLocation, domain.DeckLocation)
	}

	// Add one matching card to player's hand (Asso di Denari)
	service.gameState.Deck.MoveCardMatching(domain.DeckLocation, domain.PlayerHandLocation, domain.Asso, domain.Denari)

	// Get the player card for selection
	playerHandCards := service.gameState.Deck.CardsAt(domain.PlayerHandLocation)
	if len(playerHandCards) != 1 {
		t.Fatalf("Expected 1 card in player's hand, got %d", len(playerHandCards))
	}
	playerCard := playerHandCards[0]

	// Verify table has one card and player hand has one card
	tableCount := len(service.gameState.Deck.CardsAt(domain.TableLocation))
	playerHandCount := len(service.gameState.Deck.CardsAt(domain.PlayerHandLocation))
	deckCount := len(service.gameState.Deck.CardsAt(domain.DeckLocation))

	t.Logf("Before capture: table=%d, playerHand=%d, deck=%d", tableCount, playerHandCount, deckCount)
	assert.Equal(t, 1, tableCount, "Table should have one card")
	assert.Equal(t, 1, playerHandCount, "Player should have one card")

	// Make sure it's player's turn
	service.gameState.Status = domain.StatusPlayerTurn

	// Ensure scopa count starts at 0
	service.playerScopaCount = 0

	// Select the player card
	service.selectedCard = playerCard
	t.Logf("Selected card: %v", playerCard)

	// Verify we have capture options
	captureOptions := service.findCaptureOptions(playerCard)
	t.Logf("Capture options: %v", captureOptions)
	assert.Equal(t, 1, len(captureOptions), "Should have one capture option")

	// Perform capture using direct table card reference
	captureCards := service.gameState.Deck.CardsAt(domain.TableLocation)
	t.Logf("Capture cards: %v", captureCards)
	service.CaptureCombination(captureCards)

	// Check results after capture
	tablePostCount := len(service.gameState.Deck.CardsAt(domain.TableLocation))
	playerHandPostCount := len(service.gameState.Deck.CardsAt(domain.PlayerHandLocation))
	captureCount := len(service.gameState.Deck.CardsAt(domain.PlayerCapturesLocation))

	t.Logf("After capture: table=%d, playerHand=%d, captures=%d, scopa=%d",
		tablePostCount, playerHandPostCount, captureCount, service.playerScopaCount)

	// Check that the player got a scopa point
	assert.Equal(t, 1, service.playerScopaCount, "Player should get a scopa point")

	// Check UI model
	model := service.GetUIModel()
	assert.Equal(t, 1, model.Score.Components[4].PlayerScore, "UI should show player's scopa point")
	assert.Equal(t, 1, model.Score.Components[4].PlayerCardCount, "UI should show player's scopa count")
}

func TestAIScopa(t *testing.T) {
	// Setup
	service := NewGameService()
	service.StartNewGame()

	// Set up for a clean test
	service.SetupCombinationTest() // This creates a fresh game with specific cards

	// Clear the table
	tableCards := service.gameState.Deck.CardsAt(domain.TableLocation)
	for _, card := range tableCards {
		service.gameState.Deck.MoveCard(card, domain.TableLocation, domain.DeckLocation)
	}

	// Add one card to the table (Asso di Coppe)
	service.gameState.Deck.MoveCardMatching(domain.DeckLocation, domain.TableLocation, domain.Asso, domain.Coppe)

	// Clear AI hand
	aiCards := service.gameState.Deck.CardsAt(domain.AIHandLocation)
	for _, card := range aiCards {
		service.gameState.Deck.MoveCard(card, domain.AIHandLocation, domain.DeckLocation)
	}

	// Add one matching card to AI's hand (Asso di Denari)
	service.gameState.Deck.MoveCardMatching(domain.DeckLocation, domain.AIHandLocation, domain.Asso, domain.Denari)

	// Set it to AI's turn
	service.gameState.Status = domain.StatusAITurn

	// Ensure scopa count starts at 0
	service.aiScopaCount = 0

	// Log the state before AI turn
	tableCount := len(service.gameState.Deck.CardsAt(domain.TableLocation))
	aiHandCount := len(service.gameState.Deck.CardsAt(domain.AIHandLocation))
	deckCount := len(service.gameState.Deck.CardsAt(domain.DeckLocation))

	t.Logf("Before AI turn: table=%d, aiHand=%d, deck=%d", tableCount, aiHandCount, deckCount)
	assert.Equal(t, 1, tableCount, "Table should have one card")
	assert.Equal(t, 1, aiHandCount, "AI should have one card")

	// Play AI turn
	service.PlayAITurn()

	// Check results after AI turn
	tablePostCount := len(service.gameState.Deck.CardsAt(domain.TableLocation))
	aiHandPostCount := len(service.gameState.Deck.CardsAt(domain.AIHandLocation))
	captureCount := len(service.gameState.Deck.CardsAt(domain.AICapturesLocation))

	t.Logf("After AI turn: table=%d, aiHand=%d, captures=%d, scopa=%d",
		tablePostCount, aiHandPostCount, captureCount, service.aiScopaCount)

	// Check that the AI got a scopa point
	assert.Equal(t, 1, service.aiScopaCount, "AI should get a scopa point")

	// Check UI model
	model := service.GetUIModel()
	assert.Equal(t, 1, model.Score.Components[4].AIScore, "UI should show AI's scopa point")
	assert.Equal(t, 1, model.Score.Components[4].AICardCount, "UI should show AI's scopa count")
}

func TestNoScopaForLastCard(t *testing.T) {
	// Setup
	service := NewGameService()
	service.StartNewGame()

	// Use our test combination setup for a clean test environment
	service.SetupCombinationTest()

	// Add one card to the table with Asso rank
	service.gameState.Deck.MoveCardMatching(domain.DeckLocation, domain.TableLocation, domain.Asso, domain.Coppe)

	// Add one matching Asso to player's hand
	service.gameState.Deck.MoveCardMatching(domain.DeckLocation, domain.PlayerHandLocation, domain.Asso, domain.Denari)

	// Empty the deck completely by moving all cards to player's captures
	deckCards := service.gameState.Deck.CardsAt(domain.DeckLocation)
	for _, card := range deckCards {
		service.gameState.Deck.MoveCard(card, domain.DeckLocation, domain.PlayerCapturesLocation)
	}

	// Get the player card for selection
	playerHandCards := service.gameState.Deck.CardsAt(domain.PlayerHandLocation)
	if len(playerHandCards) != 1 {
		t.Fatalf("Expected 1 card in player's hand, got %d", len(playerHandCards))
	}
	playerCard := playerHandCards[0]

	// Select the player card
	service.selectedCard = playerCard

	// Verify initial state
	tableCardsCount := len(service.gameState.Deck.CardsAt(domain.TableLocation))
	playerHandCount := len(service.gameState.Deck.CardsAt(domain.PlayerHandLocation))
	deckCardsCount := len(service.gameState.Deck.CardsAt(domain.DeckLocation))
	aiHandCount := len(service.gameState.Deck.CardsAt(domain.AIHandLocation))

	t.Logf("Initial state: table=%d, playerHand=%d, deck=%d, aiHand=%d",
		tableCardsCount, playerHandCount, deckCardsCount, aiHandCount)

	assert.Equal(t, 1, tableCardsCount, "Table should have 1 card")
	assert.Equal(t, 1, playerHandCount, "Player should have 1 card")
	assert.Equal(t, 0, deckCardsCount, "Deck should be empty")
	assert.Equal(t, 0, aiHandCount, "AI hand should be empty")

	// Perform capture using direct table card reference
	captureCards := service.gameState.Deck.CardsAt(domain.TableLocation)
	service.CaptureCombination(captureCards)

	// Log state after capture
	t.Logf("After capture: playerHand=%d, deck=%d",
		len(service.gameState.Deck.CardsAt(domain.PlayerHandLocation)),
		len(service.gameState.Deck.CardsAt(domain.DeckLocation)))
	t.Logf("Player scopa count: %d", service.playerScopaCount)

	// Check that the player did NOT get a scopa point (last card exception)
	assert.Equal(t, 0, service.playerScopaCount, "No scopa should be awarded for last card")

	// Check UI model
	model := service.GetUIModel()
	assert.Equal(t, 0, model.Score.Components[4].PlayerScore, "UI should show no scopa points")
	assert.Equal(t, 0, model.Score.Components[4].PlayerCardCount, "UI should show no scopa count")
}
