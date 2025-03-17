package views

import (
	"fmt"
	"strings"

	"github.com/xpmatteo/scopa-trainer/pkg/domain"
)

// CardImagePath returns the path to the image file for a given card
func CardImagePath(card domain.Card) string {
	// Map suit to its single-letter code
	var suitCode string
	switch card.Suit {
	case domain.Spade:
		suitCode = "s"
	case domain.Coppe:
		suitCode = "c"
	case domain.Bastoni:
		suitCode = "b"
	case domain.Denari:
		suitCode = "d"
	default:
		suitCode = "s" // Default to spade if unknown
	}

	// Format the path
	return fmt.Sprintf("/static/images/cards/Napoletane/%d%s.jpg", card.Rank, suitCode)
}

// CardBackImagePath returns the path to the card back image
func CardBackImagePath() string {
	return "/static/images/cards/Napoletane/bg.jpg"
}

// SuitToLower converts a suit name to lowercase
func SuitToLower(suit domain.Suit) string {
	return strings.ToLower(string(suit))
}
