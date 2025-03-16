package domain

// GameStatus represents the current state of the game
type GameStatus int

const (
	GameNotStarted GameStatus = iota
	PlayerTurn
	AITurn
)

// String returns a string representation of the GameStatus
func (s GameStatus) String() string {
	switch s {
	case GameNotStarted:
		return "Game Not Started"
	case PlayerTurn:
		return "Player's Turn"
	case AITurn:
		return "AI's Turn"
	default:
		return "Unknown Status"
	}
}
