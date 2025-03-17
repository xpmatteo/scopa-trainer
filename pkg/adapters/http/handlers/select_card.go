package handlers

import (
	"html/template"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/xpmatteo/scopa-trainer/pkg/domain"
)

// UIModelProvider defines the interface for getting the UI model
type UIModelProvider interface {
	GetUIModel() domain.UIModel
}

// NewHandleIndex creates a handler for the index page
func NewHandleIndex(provider UIModelProvider, templ *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		model := provider.GetUIModel()
		if err := templ.Execute(w, model); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// GameStarter defines the interface for starting a new game
type GameStarter interface {
	StartNewGame()
}

// NewHandleNewGame creates a handler for starting a new game
func NewHandleNewGame(starter GameStarter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		starter.StartNewGame()
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// CardSelector defines the interface for selecting a card
type CardSelector interface {
	SelectCard(suit domain.Suit, rank domain.Rank)
}

// NewHandleSelectCard creates a handler for selecting a card
func NewHandleSelectCard(selector CardSelector) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse form values
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		suit := domain.Suit(r.PostForm.Get("suit"))
		rankStr := r.PostForm.Get("rank")
		rankInt, err := strconv.Atoi(rankStr)
		if err != nil {
			http.Error(w, "Invalid rank parameter", http.StatusBadRequest)
			return
		}
		rank := domain.Rank(rankInt)

		// Process the action
		selector.SelectCard(suit, rank)

		// Redirect to the main page
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// SelectedCardPlayer defines the interface for playing a selected card
type SelectedCardPlayer interface {
	PlaySelectedCard()
}

// NewHandlePlayCard creates a handler for playing a selected card
func NewHandlePlayCard(player SelectedCardPlayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		player.PlaySelectedCard()

		// Redirect to the main page
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// AITurnPlayer defines the interface for playing the AI's turn
type AITurnPlayer interface {
	PlayAITurn()
}

// NewHandleAITurn creates a handler for triggering the AI's turn
func NewHandleAITurn(player AITurnPlayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		player.PlayAITurn()

		// Redirect to the main page
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

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

// NewHandleReviewGame creates a handler for reviewing the game
// This is a placeholder for future functionality
func NewHandleReviewGame() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// For now, just redirect back to the game page
		// This is a placeholder for future functionality
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
