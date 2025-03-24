package domain

import (
	"fmt"
	"math/rand"
	"time"
)

// Location represents where a card can be in the game
type Location string

// Location constants
const (
	DeckLocation           Location = "deck"
	PlayerHandLocation     Location = "player_hand"
	AIHandLocation         Location = "ai_hand"
	TableLocation          Location = "table"
	PlayerCapturesLocation Location = "player_captures"
	AICapturesLocation     Location = "ai_captures"
)

// Deck represents a mapping of cards to their locations
type Deck struct {
	cardLocations map[Card]Location
	orderedCards  []Card // Maintains order for deck operations
}

// NewDeck creates a new deck with all cards in the deck location
func NewDeck() *Deck {
	deck := &Deck{
		cardLocations: make(map[Card]Location),
		orderedCards:  make([]Card, 0, 40),
	}

	// Create all 40 cards and put them in the deck location
	for _, suit := range AllSuits() {
		for _, rank := range AllRanks() {
			card := Card{
				Suit: suit,
				Rank: rank,
			}
			deck.cardLocations[card] = DeckLocation
			deck.orderedCards = append(deck.orderedCards, card)
		}
	}

	return deck
}

// CardsAt returns all cards at a given location
func (d *Deck) CardsAt(location Location) []Card {
	var result []Card
	for _, card := range d.orderedCards {
		if d.cardLocations[card] == location {
			result = append(result, card)
		}
	}
	return result
}

// DealCards moves n cards from source to destination
func (d *Deck) DealCards(source, destination Location, n int) {
	sourceCards := d.CardsAt(source)
	if n > len(sourceCards) {
		n = len(sourceCards)
	}

	// Update each card's location
	for i := 0; i < n; i++ {
		d.cardLocations[sourceCards[i]] = destination
	}
}

// MoveCard moves a card from source to destination
func (d *Deck) MoveCard(card Card, source, destination Location) {
	if d.cardLocations[card] == source {
		d.cardLocations[card] = destination
	}
}

// GetCardLocation returns the location of a card
func (d *Deck) GetCardLocation(card Card) Location {
	return d.cardLocations[card]
}

// Shuffle randomizes the order of cards in the deck
func (d *Deck) Shuffle() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Get cards in deck location
	deckCards := d.CardsAt(DeckLocation)

	// Shuffle only those cards
	r.Shuffle(len(deckCards), func(i, j int) {
		deckCards[i], deckCards[j] = deckCards[j], deckCards[i]
	})

	// Update the order in orderedCards
	var newOrder []Card

	// Add non-deck cards first (maintaining their original order)
	for _, card := range d.orderedCards {
		if d.cardLocations[card] != DeckLocation {
			newOrder = append(newOrder, card)
		}
	}

	// Add shuffled deck cards
	newOrder = append(newOrder, deckCards...)

	// Replace the ordered cards
	d.orderedCards = newOrder
}

func (d *Deck) MoveCardMatching(locationFrom Location, locationTo Location, rank Rank, suit Suit) {
	for _, card := range d.CardsAt(locationFrom) {
		if card.Rank == rank && card.Suit == suit {
			d.MoveCard(card, locationFrom, locationTo)
			return
		}
	}
	panic(fmt.Sprintf("card not found: %v di %v", rank, suit))
}
