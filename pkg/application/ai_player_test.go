package application

import (
	"fmt"
	"testing"

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

	if !aiCardFound {
		t.Fatalf("Could not find any suitable AI card in the deck")
	}

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

	if !tableCardFound {
		t.Fatalf("Could not find a matching table card in the deck")
	}

	// Move the card to the table
	service.gameState.Deck.MoveCard(tableCard, domain.DeckLocation, domain.TableLocation)

	// Debug: Print the state before AI turn
	fmt.Println("Before AI turn:")
	fmt.Printf("AI hand: %v\n", service.gameState.Deck.CardsAt(domain.AIHandLocation))
	fmt.Printf("Table: %v\n", service.gameState.Deck.CardsAt(domain.TableLocation))

	// Set it to AI's turn
	service.gameState.PlayerTurn = false

	// Execute AI turn
	service.PlayAITurn()

	// Debug: Print the state after AI turn
	fmt.Println("After AI turn:")
	fmt.Printf("AI hand: %v\n", service.gameState.Deck.CardsAt(domain.AIHandLocation))
	fmt.Printf("Table: %v\n", service.gameState.Deck.CardsAt(domain.TableLocation))
	fmt.Printf("AI captures: %v\n", service.gameState.Deck.CardsAt(domain.AICapturesLocation))

	// Verify that AI captured the matching card
	aiCaptures := service.gameState.Deck.CardsAt(domain.AICapturesLocation)
	if len(aiCaptures) != 2 {
		t.Errorf("Expected AI to have 2 captured cards, got %d", len(aiCaptures))
	}

	// Verify that the table is now empty
	tableCards = service.gameState.Deck.CardsAt(domain.TableLocation)
	if len(tableCards) != 0 {
		t.Errorf("Expected table to be empty, got %d cards", len(tableCards))
	}

	// Verify that it's now the player's turn
	if !service.gameState.PlayerTurn {
		t.Errorf("Expected it to be player's turn after AI move")
	}
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
	if len(deckCards) == 0 {
		t.Fatalf("No cards in the deck")
	}

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
	if len(deckCards) == 0 {
		t.Fatalf("No cards in the deck after moving AI card")
	}

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

	if !tableCardFound {
		t.Fatalf("Could not find a card with a different rank in the deck")
	}

	// Move the card to the table
	service.gameState.Deck.MoveCard(tableCard, domain.DeckLocation, domain.TableLocation)

	// Debug: Print the state before AI turn
	fmt.Println("Before AI turn (no capture):")
	fmt.Printf("AI hand: %v\n", service.gameState.Deck.CardsAt(domain.AIHandLocation))
	fmt.Printf("Table: %v\n", service.gameState.Deck.CardsAt(domain.TableLocation))

	// Set it to AI's turn
	service.gameState.PlayerTurn = false

	// Execute AI turn
	service.PlayAITurn()

	// Debug: Print the state after AI turn
	fmt.Println("After AI turn (no capture):")
	fmt.Printf("AI hand: %v\n", service.gameState.Deck.CardsAt(domain.AIHandLocation))
	fmt.Printf("Table: %v\n", service.gameState.Deck.CardsAt(domain.TableLocation))
	fmt.Printf("AI captures: %v\n", service.gameState.Deck.CardsAt(domain.AICapturesLocation))

	// Verify that AI did not capture any card
	aiCaptures := service.gameState.Deck.CardsAt(domain.AICapturesLocation)
	if len(aiCaptures) != 0 {
		t.Errorf("Expected AI to have 0 captured cards, got %d", len(aiCaptures))
	}

	// Verify that the table now has 2 cards (the original card and the played card)
	tableCards = service.gameState.Deck.CardsAt(domain.TableLocation)
	if len(tableCards) != 2 {
		t.Errorf("Expected table to have 2 cards, got %d cards", len(tableCards))
		fmt.Printf("Table cards: %v\n", tableCards)
	}

	// Verify that it's now the player's turn
	if !service.gameState.PlayerTurn {
		t.Errorf("Expected it to be player's turn after AI move")
	}
}
