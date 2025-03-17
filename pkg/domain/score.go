package domain

// ScoreComponent represents a single scoring component in Scopa
type ScoreComponent struct {
	Name            string
	Description     string
	PlayerScore     int
	AIScore         int
	PlayerCardCount int
	AICardCount     int
}

// Score represents the complete score for a game of Scopa
type Score struct {
	Components  []ScoreComponent
	PlayerTotal int
	AITotal     int
}

// NewScore creates a new Score with initialized components
func NewScore() Score {
	return Score{
		Components: []ScoreComponent{
			{
				Name:            "Carte",
				Description:     "Most cards captured",
				PlayerScore:     0,
				AIScore:         0,
				PlayerCardCount: 0,
				AICardCount:     0,
			},
			{
				Name:            "Ori",
				Description:     "Most Denari cards captured",
				PlayerScore:     0,
				AIScore:         0,
				PlayerCardCount: 0,
				AICardCount:     0,
			},
			{
				Name:            "Primiera",
				Description:     "Most Sette cards captured",
				PlayerScore:     0,
				AIScore:         0,
				PlayerCardCount: 0,
				AICardCount:     0,
			},
			{
				Name:            "Settebello",
				Description:     "Captured the Sette di Denari",
				PlayerScore:     0,
				AIScore:         0,
				PlayerCardCount: 0,
				AICardCount:     0,
			},
		},
		PlayerTotal: 0,
		AITotal:     0,
	}
}

// CalculateTotals updates the total scores based on the components
func (s *Score) CalculateTotals() {
	s.PlayerTotal = 0
	s.AITotal = 0

	for _, component := range s.Components {
		s.PlayerTotal += component.PlayerScore
		s.AITotal += component.AIScore
	}
}

// CalculateScore calculates the score based on the captured cards
func CalculateScore(playerCards []Card, aiCards []Card) Score {
	score := NewScore()

	// Calculate "Carte" - Most cards captured
	playerCardCount := len(playerCards)
	aiCardCount := len(aiCards)

	// Set the card counts for display
	score.Components[0].PlayerCardCount = playerCardCount
	score.Components[0].AICardCount = aiCardCount

	if playerCardCount > aiCardCount {
		score.Components[0].PlayerScore = 1
	} else if aiCardCount > playerCardCount {
		score.Components[0].AIScore = 1
	}
	// If tied, no points awarded

	// Calculate "Ori" - Most Denari cards captured
	playerDenariCount := countCardsBySuit(playerCards, Denari)
	aiDenariCount := countCardsBySuit(aiCards, Denari)

	// Set the Denari counts for display
	score.Components[1].PlayerCardCount = playerDenariCount
	score.Components[1].AICardCount = aiDenariCount

	if playerDenariCount > aiDenariCount {
		score.Components[1].PlayerScore = 1
	} else if aiDenariCount > playerDenariCount {
		score.Components[1].AIScore = 1
	}
	// If tied, no points awarded

	// Calculate "Primiera" - Most Sette cards captured
	playerSetteCount := countCardsByRank(playerCards, Sette)
	aiSetteCount := countCardsByRank(aiCards, Sette)

	// Set the Sette counts for display
	score.Components[2].PlayerCardCount = playerSetteCount
	score.Components[2].AICardCount = aiSetteCount

	if playerSetteCount > aiSetteCount {
		score.Components[2].PlayerScore = 1
	} else if aiSetteCount > playerSetteCount {
		score.Components[2].AIScore = 1
	}
	// If tied, no points awarded

	// Calculate "Settebello" - Captured the Sette di Denari
	settebello := Card{Suit: Denari, Rank: Sette}

	// For Settebello, the count is either 1 or 0
	score.Components[3].PlayerCardCount = 0
	score.Components[3].AICardCount = 0

	if containsCard(playerCards, settebello) {
		score.Components[3].PlayerScore = 1
		score.Components[3].PlayerCardCount = 1
	} else if containsCard(aiCards, settebello) {
		score.Components[3].AIScore = 1
		score.Components[3].AICardCount = 1
	}

	// Calculate totals
	score.CalculateTotals()

	return score
}

// countCardsBySuit counts the number of cards with the given suit
func countCardsBySuit(cards []Card, suit Suit) int {
	count := 0
	for _, card := range cards {
		if card.Suit == suit {
			count++
		}
	}
	return count
}

// countCardsByRank counts the number of cards with the given rank
func countCardsByRank(cards []Card, rank Rank) int {
	count := 0
	for _, card := range cards {
		if card.Rank == rank {
			count++
		}
	}
	return count
}

// containsCard checks if the given card is in the slice
func containsCard(cards []Card, card Card) bool {
	for _, c := range cards {
		if c.Suit == card.Suit && c.Rank == card.Rank {
			return true
		}
	}
	return false
}
