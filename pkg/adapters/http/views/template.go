package views

import (
	"html/template"
	"strings"
)

func ParseTemplates(files ...string) *template.Template {
	funcMap := template.FuncMap{
		"lower": strings.ToLower,
	}
	templ, err := template.New("game.html").Funcs(funcMap).ParseFiles(files...)
	if err != nil {
		panic(err)
	}
	return templ
}
