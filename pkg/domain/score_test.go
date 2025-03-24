package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewScore(t *testing.T) {
	score := NewScore()

	// Verify initial state
	assert.Equal(t, 5, len(score.Components), "Should have 5 scoring components")
	assert.Equal(t, "Carte", score.Components[0].Name)
	assert.Equal(t, "Ori", score.Components[1].Name)
	assert.Equal(t, "Primiera", score.Components[2].Name)
	assert.Equal(t, "Settebello", score.Components[3].Name)
	assert.Equal(t, "Scope", score.Components[4].Name)
	assert.Equal(t, 0, score.PlayerTotal)
	assert.Equal(t, 0, score.AITotal)
}

func TestCalculateTotals(t *testing.T) {
	score := NewScore()

	// Set some scores
	score.Components[0].PlayerScore = 1 // Carte
	score.Components[1].AIScore = 1     // Ori
	score.Components[2].PlayerScore = 1 // Primiera
	score.Components[3].AIScore = 1     // Settebello
	score.Components[4].PlayerScore = 2 // Scope

	// Calculate totals
	score.CalculateTotals()

	// Verify totals
	assert.Equal(t, 4, score.PlayerTotal)
	assert.Equal(t, 2, score.AITotal)
}

func TestCalculateScore(t *testing.T) {
	tests := []struct {
		name                string
		playerCards         []Card
		aiCards             []Card
		expectedPlayerTotal int
		expectedAITotal     int
		expectedComponents  map[string]struct{ player, ai int }
	}{
		{
			name: "Player wins all categories",
			playerCards: []Card{
				{Suit: Denari, Rank: Asso},
				{Suit: Denari, Rank: Due},
				{Suit: Denari, Rank: Tre},
				{Suit: Denari, Rank: Sette},  // Settebello
				{Suit: Coppe, Rank: Sette},   // Extra Sette for Primiera
				{Suit: Bastoni, Rank: Sette}, // Extra Sette for Primiera
			},
			aiCards: []Card{
				{Suit: Spade, Rank: Asso},
				{Suit: Spade, Rank: Due},
				{Suit: Coppe, Rank: Tre},
			},
			expectedPlayerTotal: 4,
			expectedAITotal:     0,
			expectedComponents: map[string]struct{ player, ai int }{
				"Carte":      {1, 0},
				"Ori":        {1, 0},
				"Primiera":   {1, 0},
				"Settebello": {1, 0},
				"Scope":      {0, 0}, // Scope component should be initialized as 0
			},
		},
		{
			name: "AI wins all categories",
			playerCards: []Card{
				{Suit: Spade, Rank: Asso},
				{Suit: Spade, Rank: Due},
			},
			aiCards: []Card{
				{Suit: Denari, Rank: Asso},
				{Suit: Denari, Rank: Due},
				{Suit: Denari, Rank: Tre},
				{Suit: Denari, Rank: Sette},  // Settebello
				{Suit: Coppe, Rank: Sette},   // Extra Sette for Primiera
				{Suit: Bastoni, Rank: Sette}, // Extra Sette for Primiera
			},
			expectedPlayerTotal: 0,
			expectedAITotal:     4,
			expectedComponents: map[string]struct{ player, ai int }{
				"Carte":      {0, 1},
				"Ori":        {0, 1},
				"Primiera":   {0, 1},
				"Settebello": {0, 1},
				"Scope":      {0, 0}, // Scope component should be initialized as 0
			},
		},
		{
			name: "Tie in card count",
			playerCards: []Card{
				{Suit: Spade, Rank: Asso},
				{Suit: Spade, Rank: Due},
				{Suit: Denari, Rank: Sette}, // Settebello
			},
			aiCards: []Card{
				{Suit: Coppe, Rank: Asso},
				{Suit: Coppe, Rank: Due},
				{Suit: Coppe, Rank: Sette}, // Sette for Primiera
			},
			expectedPlayerTotal: 2,
			expectedAITotal:     0,
			expectedComponents: map[string]struct{ player, ai int }{
				"Carte":      {0, 0}, // Tie, no points
				"Ori":        {1, 0}, // Player has 1 Denari, AI has 0 Denari
				"Primiera":   {0, 0}, // Tie in Sette cards (1 each)
				"Settebello": {1, 0}, // Player has Settebello
				"Scope":      {0, 0}, // Scope component should be initialized as 0
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			score := CalculateScore(test.playerCards, test.aiCards)

			// Check totals
			assert.Equal(t, test.expectedPlayerTotal, score.PlayerTotal, "Player total score should match")
			assert.Equal(t, test.expectedAITotal, score.AITotal, "AI total score should match")

			// Check individual components
			for name, expected := range test.expectedComponents {
				for _, component := range score.Components {
					if component.Name == name {
						assert.Equal(t, expected.player, component.PlayerScore, "Player score for %s should match", name)
						assert.Equal(t, expected.ai, component.AIScore, "AI score for %s should match", name)
					}
				}
			}
		})
	}
}
