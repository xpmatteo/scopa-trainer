package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGameService(t *testing.T) {
	service := NewGameService()
	model := service.GetUIModel()
	assert.Equal(t, "hello player", model.GamePrompt)
}
