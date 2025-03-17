package domain

// GameStatus represents the current state of the game
type GameStatus int

const (
	StatusGameNotStarted GameStatus = iota
	StatusPlayerTurn
	StatusAITurn
	StatusGameOver
)

// String returns a string representation of the GameStatus
func (s GameStatus) String() string {
	switch s {
	case StatusGameNotStarted:
		return "Game Not Started"
	case StatusPlayerTurn:
		return "Player's Turn"
	case StatusAITurn:
		return "AI's Turn"
	case StatusGameOver:
		return "Game Over"
	default:
		panic("Unknown Status")
	}
}
