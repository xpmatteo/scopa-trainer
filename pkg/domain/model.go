package domain

// UIModel represents the complete state of the UI at any point in time
type UIModel struct {
	GamePrompt string
}

// NewUIModel creates a new UI model with initial state
func NewUIModel() UIModel {
	return UIModel{
		GamePrompt: "hello player",
	}
}
