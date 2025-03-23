package views

import (
	"html/template"
	"strings"

	"github.com/xpmatteo/scopa-trainer/pkg/domain"
)

// CardInCaptureOptions checks if a card is part of any capture option
func CardInCaptureOptions(card domain.Card, captureOptions [][]domain.Card) bool {
	for _, combo := range captureOptions {
		for _, c := range combo {
			if c.Suit == card.Suit && c.Rank == card.Rank {
				return true
			}
		}
	}
	return false
}

// ParseTemplates parses the template files and adds custom functions
func ParseTemplates(files ...string) *template.Template {
	funcMap := template.FuncMap{
		"lower":                strings.ToLower,
		"cardImagePath":        CardImagePath,
		"cardBackImagePath":    CardBackImagePath,
		"suitToLower":          SuitToLower,
		"cardInCaptureOptions": CardInCaptureOptions,
	}
	templ, err := template.New("game.html").Funcs(funcMap).ParseFiles(files...)
	if err != nil {
		panic(err)
	}
	return templ
}
