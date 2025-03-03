package domain

// Card represents a playing card in the Italian deck
type Card struct {
	Suit  Suit
	Rank  Rank
	Name  string
	Value int
}

// GameState represents the current state of the game
type GameState struct {
	Deck       *Deck
	PlayerTurn bool
}

// UIModel represents the complete state of the UI at any point in time
type UIModel struct {
	GamePrompt        string
	ShowNewGameButton bool
	TableCards        []Card
	PlayerHand        []Card
	GameInProgress    bool
	PlayerTurn        bool
}

// NewUIModel creates a new UI model with initial state
func NewUIModel() UIModel {
	return UIModel{
		GamePrompt:        "Welcome to Scopa Trainer! Click 'New Game' to start playing.",
		ShowNewGameButton: true,
		TableCards:        []Card{},
		PlayerHand:        []Card{},
		GameInProgress:    false,
		PlayerTurn:        false,
	}
}

// NewGameState initializes a new game state
func NewGameState() GameState {
	deck := NewDeck()
	deck.Shuffle()

	// Deal cards to each player
	deck.DealCards(DeckLocation, PlayerHandLocation, 10)
	deck.DealCards(DeckLocation, AIHandLocation, 10)

	return GameState{
		Deck:       deck,
		PlayerTurn: true, // Player goes first by default
	}
}
