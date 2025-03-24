package domain

// UIModel represents the complete state of the UI at any point in time
type UIModel struct {
	GamePrompt          string
	ShowNewGameButton   bool
	TableCards          []Card
	PlayerHand          []Card
	GameInProgress      bool
	PlayerTurn          bool
	SelectedCard        Card
	CanPlaySelectedCard bool
	DeckCount           int
	PlayerCaptureCount  int
	AICaptureCount      int
	GameOver            bool
	PlayerCaptureCards  []Card
	AICaptureCards      []Card
	Score               Score
	CaptureOptions      [][]Card // All possible capture combinations for selected card
	SelectedTableCards  []Card   // Currently selected table cards for combination capture
	CanConfirmCapture   bool     // Whether player can confirm the current capture selection
	LastAICardPlayed    Card     // The card the AI played in its last move
	LastAICapture       []Card   // The cards the AI captured in its last move
	ShowAIMove          bool     // Whether to show the AI's last move
}

var NO_CARD_SELECTED = Card{}

// NewUIModel creates a new UI model with initial state
func NewUIModel() UIModel {
	return UIModel{
		GamePrompt:          "Welcome to Scopa Trainer! Click 'New Game' to start playing.",
		ShowNewGameButton:   true,
		TableCards:          []Card{},
		PlayerHand:          []Card{},
		GameInProgress:      false,
		PlayerTurn:          false,
		SelectedCard:        NO_CARD_SELECTED,
		CanPlaySelectedCard: false,
		DeckCount:           0,
		PlayerCaptureCount:  0,
		AICaptureCount:      0,
		GameOver:            false,
		PlayerCaptureCards:  []Card{},
		AICaptureCards:      []Card{},
		Score:               NewScore(),
		CaptureOptions:      [][]Card{},
		SelectedTableCards:  []Card{},
		CanConfirmCapture:   false,
		LastAICardPlayed:    NO_CARD_SELECTED,
		LastAICapture:       []Card{},
		ShowAIMove:          false,
	}
}
