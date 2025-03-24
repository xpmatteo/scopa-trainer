package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xpmatteo/scopa-trainer/pkg/domain"
)

func TestRandomAIPlayer(t *testing.T) {
	// Create a new game service
	service := NewGameService()
	service.StartNewGame()

	// Set up a known game state for testing
	// We'll manually set up the deck to have known cards in specific locations

	// Clear existing cards from AI hand and add a specific card
	aiCards := service.gameState.Deck.CardsAt(domain.AIHandLocation)
	for _, card := range aiCards {
		service.gameState.Deck.MoveCard(card, domain.AIHandLocation, domain.DeckLocation)
	}

	// Try different suits and ranks until we find a card for the AI hand
	possibleSuits := []domain.Suit{domain.Coppe, domain.Denari, domain.Bastoni, domain.Spade}
	possibleRanks := []domain.Rank{domain.Sette, domain.Sei, domain.Cinque, domain.Quattro, domain.Tre}

	var aiCard domain.Card
	aiCardFound := false

	for _, suit := range possibleSuits {
		for _, rank := range possibleRanks {
			candidateCard := domain.Card{Suit: suit, Rank: rank}

			// Check if this card exists in the deck
			deckCards := service.gameState.Deck.CardsAt(domain.DeckLocation)
			for _, card := range deckCards {
				if card.Suit == candidateCard.Suit && card.Rank == candidateCard.Rank {
					aiCardFound = true
					aiCard = card // Use the actual card from the deck
					break
				}
			}

			if aiCardFound {
				break
			}
		}

		if aiCardFound {
			break
		}
	}

	assert.True(t, aiCardFound, "Should be able to find a suitable AI card in the deck")

	// Move the card to AI hand
	service.gameState.Deck.MoveCard(aiCard, domain.DeckLocation, domain.AIHandLocation)

	// Clear table and add a card with the same rank to the table
	tableCards := service.gameState.Deck.CardsAt(domain.TableLocation)
	for _, card := range tableCards {
		service.gameState.Deck.MoveCard(card, domain.TableLocation, domain.DeckLocation)
	}

	// Find a card with the same rank but different suit for the table
	var tableCard domain.Card
	tableCardFound := false

	for _, suit := range possibleSuits {
		// Skip the suit of the AI card
		if suit == aiCard.Suit {
			continue
		}

		candidateCard := domain.Card{Suit: suit, Rank: aiCard.Rank}

		// Check if this card exists in the deck
		deckCards := service.gameState.Deck.CardsAt(domain.DeckLocation)
		for _, card := range deckCards {
			if card.Suit == candidateCard.Suit && card.Rank == candidateCard.Rank {
				tableCardFound = true
				tableCard = card // Use the actual card from the deck
				break
			}
		}

		if tableCardFound {
			break
		}
	}

	assert.True(t, tableCardFound, "Should be able to find a matching table card in the deck")

	// Move the card to the table
	service.gameState.Deck.MoveCard(tableCard, domain.DeckLocation, domain.TableLocation)

	// Verify initial state
	aiHandBefore := service.gameState.Deck.CardsAt(domain.AIHandLocation)
	tableBefore := service.gameState.Deck.CardsAt(domain.TableLocation)
	assert.Equal(t, 1, len(aiHandBefore), "AI should have 1 card before its turn")
	assert.Equal(t, 1, len(tableBefore), "Table should have 1 card before AI turn")
	assert.Equal(t, aiCard.Rank, tableBefore[0].Rank, "Table card should have same rank as AI card")

	// Set it to AI's turn
	service.gameState.Status = domain.StatusAITurn

	// Execute AI turn
	service.PlayAITurn()

	// Verify final state
	aiHandAfter := service.gameState.Deck.CardsAt(domain.AIHandLocation)
	tableAfter := service.gameState.Deck.CardsAt(domain.TableLocation)
	aiCaptures := service.gameState.Deck.CardsAt(domain.AICapturesLocation)

	// Verify that AI captured the matching card
	assert.Equal(t, 2, len(aiCaptures), "AI should have captured 2 cards")
	assert.Equal(t, 0, len(tableAfter), "Table should be empty after AI turn")
	assert.Equal(t, 0, len(aiHandAfter), "AI hand should be empty after playing its only card")

	// Verify that it's now the player's turn
	assert.Equal(t, domain.StatusPlayerTurn, service.gameState.Status, "It should be player's turn after AI move")
}

func TestRandomAIPlayerNoCapture(t *testing.T) {
	// Create a new game service
	service := NewGameService()
	service.StartNewGame()

	// Set up a known game state for testing
	// Clear existing cards from AI hand and add a specific card
	aiCards := service.gameState.Deck.CardsAt(domain.AIHandLocation)
	for _, card := range aiCards {
		service.gameState.Deck.MoveCard(card, domain.AIHandLocation, domain.DeckLocation)
	}

	// Get all cards from the deck
	deckCards := service.gameState.Deck.CardsAt(domain.DeckLocation)
	assert.NotEmpty(t, deckCards, "Deck should have cards")

	// Use the first card from the deck for the AI hand
	aiCard := deckCards[0]
	service.gameState.Deck.MoveCard(aiCard, domain.DeckLocation, domain.AIHandLocation)

	// Clear table and add a different card to the table
	tableCards := service.gameState.Deck.CardsAt(domain.TableLocation)
	for _, card := range tableCards {
		service.gameState.Deck.MoveCard(card, domain.TableLocation, domain.DeckLocation)
	}

	// Get updated deck cards after moving the AI card
	deckCards = service.gameState.Deck.CardsAt(domain.DeckLocation)
	assert.NotEmpty(t, deckCards, "Deck should have cards after moving AI card")

	// Find a card with a different rank for the table
	var tableCard domain.Card
	tableCardFound := false

	for _, card := range deckCards {
		if card.Rank != aiCard.Rank {
			tableCard = card
			tableCardFound = true
			break
		}
	}

	assert.True(t, tableCardFound, "Should be able to find a card with a different rank in the deck")

	// Move the card to the table
	service.gameState.Deck.MoveCard(tableCard, domain.DeckLocation, domain.TableLocation)

	// Verify initial state
	aiHandBefore := service.gameState.Deck.CardsAt(domain.AIHandLocation)
	tableBefore := service.gameState.Deck.CardsAt(domain.TableLocation)
	assert.Equal(t, 1, len(aiHandBefore), "AI should have 1 card before its turn")
	assert.Equal(t, 1, len(tableBefore), "Table should have 1 card before AI turn")
	assert.NotEqual(t, aiCard.Rank, tableBefore[0].Rank, "Table card should have different rank than AI card")

	// Set it to AI's turn
	service.gameState.Status = domain.StatusAITurn

	// Execute AI turn
	service.PlayAITurn()

	// Verify final state
	aiHandAfter := service.gameState.Deck.CardsAt(domain.AIHandLocation)
	tableAfter := service.gameState.Deck.CardsAt(domain.TableLocation)
	aiCaptures := service.gameState.Deck.CardsAt(domain.AICapturesLocation)

	// Verify that AI did not capture any card
	assert.Equal(t, 0, len(aiCaptures), "AI should have 0 captured cards")
	assert.Equal(t, 2, len(tableAfter), "Table should have 2 cards after AI turn")
	assert.Equal(t, 0, len(aiHandAfter), "AI hand should be empty after playing its only card")

	// Verify that it's now the player's turn
	assert.Equal(t, domain.StatusPlayerTurn, service.gameState.Status, "It should be player's turn after AI move")
}

func TestAILastMoveTracking(t *testing.T) {
	// Create a new game service
	service := NewGameService()
	service.StartNewGame()

	// Set up a known game state for testing
	// Clear existing cards from AI hand and add a specific card
	aiCards := service.gameState.Deck.CardsAt(domain.AIHandLocation)
	for _, card := range aiCards {
		service.gameState.Deck.MoveCard(card, domain.AIHandLocation, domain.DeckLocation)
	}

	// Add a specific card to AI hand
	aiCard := domain.Card{Suit: domain.Coppe, Rank: domain.Sette}
	addCardToAIHand(service, aiCard)

	// Clear the table
	tableCards := service.gameState.Deck.CardsAt(domain.TableLocation)
	for _, card := range tableCards {
		service.gameState.Deck.MoveCard(card, domain.TableLocation, domain.DeckLocation)
	}

	// Add a matching card to the table
	tableCard := domain.Card{Suit: domain.Denari, Rank: domain.Sette}
	addCardToTable(service, tableCard)

	// Set it to AI's turn
	service.gameState.Status = domain.StatusAITurn

	// Verify initial state - no AI move should be shown yet
	initialModel := service.GetUIModel()
	assert.False(t, initialModel.ShowAIMove, "AI move should not be shown before AI plays")
	assert.Equal(t, domain.NO_CARD_SELECTED, initialModel.LastAICardPlayed, "No AI card should be tracked yet")
	assert.Empty(t, initialModel.LastAICapture, "No AI capture should be tracked yet")

	// Execute AI turn
	service.PlayAITurn()

	// Verify AI's move is tracked
	updatedModel := service.GetUIModel()
	assert.True(t, updatedModel.ShowAIMove, "AI move should be shown after AI plays")
	assert.Equal(t, aiCard, updatedModel.LastAICardPlayed, "The AI's played card should be tracked")
	assert.Equal(t, 1, len(updatedModel.LastAICapture), "AI's capture should be tracked")
	assert.Equal(t, tableCard, updatedModel.LastAICapture[0], "The captured card should match the table card")

	// Simulate player selecting a card - should clear AI move
	playerHand := service.gameState.Deck.CardsAt(domain.PlayerHandLocation)
	assert.NotEmpty(t, playerHand, "Player should have cards")
	service.SelectCard(playerHand[0].Suit, playerHand[0].Rank)

	// Verify AI move is cleared
	afterSelectionModel := service.GetUIModel()
	assert.False(t, afterSelectionModel.ShowAIMove, "AI move should be hidden after player selects a card")
}
