package handlers

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/xpmatteo/scopa-trainer/pkg/domain"
)

// GameStateManager defines the interface for testing the game over state
type GameStateManager interface {
	StartNewGame()
	SetGameOver()
	GetDeck() *domain.Deck
}

// NewHandleTestGameOver creates a handler for testing the game over state
func NewHandleTestGameOver(service GameStateManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Start a fresh game
		service.StartNewGame()

		// Get the deck
		deck := service.GetDeck()

		// Collect all cards from all locations
		allCards := getAllCards(deck)

		// Generate a random number between 0 and 20 for player cards
		rng := rand.New(rand.NewSource(time.Now().UnixNano()))
		playerCardCount := rng.Intn(21)

		// Randomly shuffle the cards
		rng.Shuffle(len(allCards), func(i, j int) {
			allCards[i], allCards[j] = allCards[j], allCards[i]
		})

		// Move cards to player and AI capture piles
		for i, card := range allCards {
			location := deck.GetCardLocation(card)

			if i < playerCardCount {
				// Move to player captures
				deck.MoveCard(card, location, domain.PlayerCapturesLocation)
			} else {
				// Move to AI captures
				deck.MoveCard(card, location, domain.AICapturesLocation)
			}
		}

		// Set the game state to game over
		service.SetGameOver()

		// Redirect to the main page to show the game over screen
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// getAllCards collects all cards from all locations
func getAllCards(deck *domain.Deck) []domain.Card {
	var allCards []domain.Card

	// Collect cards from all locations
	locations := []domain.Location{
		domain.DeckLocation,
		domain.PlayerHandLocation,
		domain.AIHandLocation,
		domain.TableLocation,
	}

	for _, location := range locations {
		allCards = append(allCards, deck.CardsAt(location)...)
	}

	return allCards
}
