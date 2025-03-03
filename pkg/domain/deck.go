package domain

import (
	"math/rand"
	"time"
)

// Suit constants
const (
	Coppe   = "Coppe"
	Denari  = "Denari"
	Bastoni = "Bastoni"
	Spade   = "Spade"
)

// Rank name mapping
var rankNames = map[int]string{
	1:  "Asso",
	2:  "Due",
	3:  "Tre",
	4:  "Quattro",
	5:  "Cinque",
	6:  "Sei",
	7:  "Sette",
	8:  "Fante",
	9:  "Cavallo",
	10: "Re",
}

// NewDeck creates a standard 40-card Italian deck
func NewDeck() []Card {
	suits := []string{Coppe, Denari, Bastoni, Spade}
	deck := make([]Card, 0, 40)

	for _, suit := range suits {
		for rank := 1; rank <= 10; rank++ {
			card := Card{
				Suit:  suit,
				Rank:  rank,
				Name:  rankNames[rank],
				Value: rank,
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
