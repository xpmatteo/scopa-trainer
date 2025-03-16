package domain

// GameState represents the current state of the game
type GameState struct {
	Deck   *Deck
	Status GameStatus
}

// NewGameState initializes a new game state
func NewGameState() GameState {
	deck := NewDeck()
	deck.Shuffle()

	// Deal cards to each player
	deck.DealCards(DeckLocation, PlayerHandLocation, 10)
	deck.DealCards(DeckLocation, AIHandLocation, 10)

	return GameState{
		Deck:   deck,
		Status: StatusPlayerTurn, // Player goes first by default
	}
}
