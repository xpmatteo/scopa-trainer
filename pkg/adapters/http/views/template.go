package views

import (
	"html/template"
	"strings"
)

// ParseTemplates parses the template files and adds custom functions
func ParseTemplates(files ...string) *template.Template {
	funcMap := template.FuncMap{
		"lower":             strings.ToLower,
		"cardImagePath":     CardImagePath,
		"cardBackImagePath": CardBackImagePath,
		"suitToLower":       SuitToLower,
	}
	templ, err := template.New("game.html").Funcs(funcMap).ParseFiles(files...)
	if err != nil {
		panic(err)
	}
	return templ
}
