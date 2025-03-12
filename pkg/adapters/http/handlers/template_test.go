package handlers

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"html/template"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xpmatteo/scopa-trainer/pkg/domain"
)

func TestWelcomeScreen(t *testing.T) {
	// Arrange
	model := domain.NewUIModel()

	// Act
	doc := renderTemplate(t, model)

	// Assert
	expected := `
		Welcome to Scopa Trainer! Click &#39;New Game&#39; to start playing. 
		New Game`
	actual := visualizeWithTestIcons(doc)
	assert.Equal(t, normalizeWhitespace(expected), actual)
}

func TestGameInProgressScreen(t *testing.T) {
	// Given a simplified game in progress model
	model := domain.NewUIModel()
	model.GameInProgress = true
	model.ShowNewGameButton = false
	model.GamePrompt = "Your turn. Select a card to play."

	// With one card on the table and one in hand
	model.TableCards = []domain.Card{
		{Suit: domain.Coppe, Rank: domain.Asso},
	}
	model.PlayerHand = []domain.Card{
		{Suit: domain.Denari, Rank: domain.Tre},
	}

	// Act
	doc := renderTemplate(t, model)

	// Assert
	expected := `
		Your turn. Select a card to play. 
		Deck: 0 cards 
		Your Captures: 0 cards 
		AI Captures: 0 cards 
		Table Cards (1) Asso Asso di Coppe 
		Your Hand (1) Tre Tre di Denari`
	actual := visualizeWithTestIcons(doc)
	assert.Equal(t, normalizeWhitespace(expected), actual)
}

func renderTemplate(t *testing.T, model domain.UIModel) string {
	funcMap := template.FuncMap{
		"lower": strings.ToLower,
	}
	templ, err := template.New("game.html").Funcs(funcMap).ParseFiles("../../../../templates/game.html")
	require.NoError(t, err)

	var buf bytes.Buffer

	if err := templ.Execute(&buf, model); err != nil {
		require.NoError(t, err)
	}
	return buf.String()
}

func visualizeWithTestIcons(html string) string {
	// Remove style element entirely
	html = replaceAll(html, `(?s)<style>.*?</style>`, "")

	// Remove script element entirely
	html = replaceAll(html, `(?s)<script>.*?</script>`, "")

	// Remove title element entirely
	html = replaceAll(html, `(?s)<title>.*?</title>`, "")

	// Replace elements with data-test-icon with their icon
	html = replaceAll(html, `<[^<>]+\bdata-test-icon="(.*?)".*?>`, " $1 ")

	// Process elements with onclick but no test-icon
	html = replaceAll(html, `onclick`, " orpo ")
	html = replaceAll(html, `\bonclick="([^"]*)"[^>]*>(.*?)`, " [👆 $1: $2] ")
	html = replaceAll(html, `<button[^>]*\bonclick="([^"]*)"[^>]*>(.*?)</button>`, " [👆 $1: $2] ")

	// Strip remaining HTML tags
	html = replaceAll(html, `</?(a|abbr|b|big|cite|code|em|i|small|span|strong|tt)\b.*?>`, "")
	html = replaceAll(html, `<[^>]*>`, " ")

	// Replace HTML entities
	html = replaceAll(html, `&nbsp;`, " ")
	html = replaceAll(html, `&lt;`, "<")
	html = replaceAll(html, `&gt;`, ">")
	html = replaceAll(html, `&quot;`, "\"")
	html = replaceAll(html, `&apos;`, "'")
	html = replaceAll(html, `&amp;`, "&")

	return normalizeWhitespace(html)
}

func replaceAll(src, regex, repl string) string {
	re := regexp.MustCompile(regex)
	return re.ReplaceAllString(src, repl)
}

func normalizeWhitespace(s string) string {
	return strings.TrimSpace(replaceAll(s, `\s+`, " "))
}
