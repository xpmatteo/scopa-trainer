package application

import (
	"github.com/xpmatteo/scopa-trainer/pkg/domain"
)

// GameService handles the application logic for the game
type GameService struct {
	model domain.UIModel
}

// NewGameService creates a new game service with initial state
func NewGameService() *GameService {
	return &GameService{
		model: domain.NewUIModel(),
	}
}

// GetUIModel returns the current UI model
func (s *GameService) GetUIModel() domain.UIModel {
	return s.model
}
