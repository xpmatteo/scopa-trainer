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
	Deck       []Card
	TableCards []Card
	PlayerHand []Card
	AIHand     []Card
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
	shuffledDeck := ShuffleDeck(deck)

	// Deal cards to each player
	playerHand, remainingDeck := DealCards(shuffledDeck, 10)
	aiHand, remainingDeck := DealCards(remainingDeck, 10)

	return GameState{
		Deck:       remainingDeck,
		TableCards: []Card{},
		PlayerHand: playerHand,
		AIHand:     aiHand,
		PlayerTurn: true, // Player goes first by default
	}
}
