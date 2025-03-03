package domain

import (
	"math/rand"
	"time"
)

// NewDeck creates a standard 40-card Italian deck
func NewDeck() []Card {
	suits := []Suit{Coppe, Denari, Bastoni, Spade}
	deck := make([]Card, 0, 40)

	for _, suit := range suits {
		for rank := Rank(1); rank <= 10; rank++ {
			card := Card{
				Suit:  suit,
				Rank:  rank,
				Name:  rankNames[rank],
				Value: int(rank),
			}
			deck = append(deck, card)
		}
	}

	return deck
}

// ShuffleDeck randomizes the order of cards in the deck
func ShuffleDeck(deck []Card) []Card {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	shuffled := make([]Card, len(deck))
	copy(shuffled, deck)

	r.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	return shuffled
}

// DealCards deals n cards from the deck and returns the dealt cards and the remaining deck
func DealCards(deck []Card, n int) ([]Card, []Card) {
	if n > len(deck) {
		n = len(deck)
	}

	dealt := deck[:n]
	remaining := deck[n:]

	return dealt, remaining
}
