package handlers

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/html"
	"html/template"
	"regexp"
	"slices"
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
		Welcome to Scopa Trainer! Click 'New Game' to start playing. 
		[👆 New Game]`
	actual := visualizeTemplate(doc)
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
		Table Cards (1) [👆 Asso di Coppe] Asso Asso di Coppe 
		Your Hand (1) [👆 Tre di Denari] Tre Tre di Denari		
		`
	actual := visualizeTemplate(doc)
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

func replaceAll(src, regex, repl string) string {
	re := regexp.MustCompile(regex)
	return re.ReplaceAllString(src, repl)
}

func normalizeWhitespace(s string) string {
	return strings.TrimSpace(replaceAll(s, `\s+`, " "))
}

func visualizeTemplate(htmlContent string) string {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return fmt.Sprintf("Error parsing HTML: %v", err)
	}

	var output strings.Builder
	skipElements := []string{"style", "script", "title"}
	visualizeNode(doc, &output, 0, skipElements)

	return normalizeWhitespace(output.String())
}

func visualizeNode(n *html.Node, output *strings.Builder, depth int, skipElements []string) {
	// Skip elements in the skip list
	if n.Type == html.ElementNode && slices.Contains(skipElements, n.Data) {
		return
	}

	if n.Type == html.ElementNode && n.Data == "button" {
		output.WriteString(fmt.Sprintf("[👆 %s] ", extractTextContent(n)))
		return
	}

	if n.Type == html.ElementNode {
		var testIcon string
		var hasOnClick bool
		var onClickValue string

		for _, attr := range n.Attr {
			if attr.Key == "data-test-icon" {
				testIcon = attr.Val
			}
			if attr.Key == "onclick" {
				hasOnClick = true
				onClickValue = attr.Val
				_ = onClickValue
			}
		}

		if testIcon != "" && hasOnClick {
			output.WriteString(fmt.Sprintf("[👆 %s] ", testIcon))
		} else if testIcon != "" {
			output.WriteString(testIcon)
			output.WriteString(" ")
		} else if hasOnClick {
			textContent := extractTextContent(n)
			output.WriteString(fmt.Sprintf("[👆 %s] ", textContent))
		}
	} else if n.Type == html.TextNode {
		text := strings.TrimSpace(n.Data)
		if text != "" {
			output.WriteString(text + " ")
		}
	}

	// Process children
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		visualizeNode(c, output, depth+1, skipElements)
	}
}

func extractTextContent(n *html.Node) string {
	var text strings.Builder
	extractTextHelper(n, &text)
	return strings.TrimSpace(text.String())
}

func extractTextHelper(n *html.Node, text *strings.Builder) {
	if n.Type == html.TextNode {
		text.WriteString(n.Data)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		extractTextHelper(c, text)
	}
}
