package handlers

import (
	"bytes"
	"html/template"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xpmatteo/scopa-trainer/pkg/domain"
	"golang.org/x/net/html"
)

// Helper function to render template and parse HTML
func renderTemplate(t *testing.T, model domain.UIModel) *html.Node {
	funcMap := template.FuncMap{
		"lower": strings.ToLower,
	}
	tmpl, err := template.New("game.html").Funcs(funcMap).ParseFiles("../../../../templates/game.html")
	assert.NoError(t, err)

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, model)
	assert.NoError(t, err)

	doc, err := html.Parse(strings.NewReader(buf.String()))
	assert.NoError(t, err)

	return doc
}

// Helper to find elements by CSS selector-like query
func findElement(node *html.Node, elementType, class string) *html.Node {
	if node == nil {
		return nil
	}

	if node.Type == html.ElementNode && node.Data == elementType {
		for _, attr := range node.Attr {
			if attr.Key == "class" && strings.Contains(attr.Val, class) {
				return node
			}
		}
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if result := findElement(c, elementType, class); result != nil {
			return result
		}
	}

	return nil
}

// Helper to find all elements matching criteria
func findAllElements(node *html.Node, elementType, class string) []*html.Node {
	var results []*html.Node

	if node == nil {
		return results
	}

	if node.Type == html.ElementNode && node.Data == elementType {
		for _, attr := range node.Attr {
			if attr.Key == "class" && strings.Contains(attr.Val, class) {
				results = append(results, node)
				break
			}
		}
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		results = append(results, findAllElements(c, elementType, class)...)
	}

	return results
}

// Helper to get text content of node
func getTextContent(node *html.Node) string {
	if node == nil {
		return ""
	}

	if node.Type == html.TextNode {
		return node.Data
	}

	var result string
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		result += getTextContent(c)
	}

	return strings.TrimSpace(result)
}

func TestWelcomeScreen(t *testing.T) {
	// Given a new game model
	model := domain.NewUIModel()

	// When rendering the template
	doc := renderTemplate(t, model)

	// Then we should see the welcome message
	promptDiv := findElement(doc, "div", "game-prompt")
	assert.NotNil(t, promptDiv, "Game prompt div should exist")
	assert.Contains(t, getTextContent(promptDiv), "Welcome to Scopa Trainer")

	// And we should see the new game button
	button := findElement(doc, "button", "new-game-button")
	assert.NotNil(t, button, "New game button should exist")
	assert.Equal(t, "New Game", getTextContent(button))

	// And we should not see the game areas
	assert.Nil(t, findElement(doc, "div", "table-area"), "Table area should not exist yet")
	assert.Nil(t, findElement(doc, "div", "hand-area"), "Hand area should not exist yet")
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

	// When rendering the template - this should not panic
	doc := renderTemplate(t, model)

	// Then we should see the game areas
	assert.NotNil(t, findElement(doc, "div", "game-prompt"), "Game prompt should be visible")
	assert.NotNil(t, findElement(doc, "div", "table-area"), "Table area should be visible")
	assert.NotNil(t, findElement(doc, "div", "hand-area"), "Hand area should be visible")

	// And we should see at least one card
	tableCards := findAllElements(doc, "div", "card")
	assert.Greater(t, len(tableCards), 0, "At least one card should be visible")
}

func TestCardSuitStyling(t *testing.T) {
	// Given a model with cards of each suit
	model := domain.NewUIModel()
	model.GameInProgress = true

	// Add one card of each suit to the table
	model.TableCards = []domain.Card{
		{Suit: domain.Coppe, Rank: domain.Asso},
		{Suit: domain.Denari, Rank: domain.Due},
		{Suit: domain.Bastoni, Rank: domain.Tre},
		{Suit: domain.Spade, Rank: domain.Quattro},
	}

	// When rendering the template
	doc := renderTemplate(t, model)

	// Then we should see cards on the table
	tableArea := findElement(doc, "div", "table-area")
	assert.NotNil(t, tableArea, "Table area should be visible")

	// And we should see cards
	cards := findAllElements(tableArea, "div", "card")
	assert.NotEmpty(t, cards, "Cards should be visible")

	// Check that we have at least one card with a suit class
	// Note: The template might be using lowercase suit names for classes
	foundSuitClass := false
	for _, card := range cards {
		for _, attr := range card.Attr {
			if attr.Key == "class" {
				classVal := strings.ToLower(attr.Val)
				if strings.Contains(classVal, "card") {
					foundSuitClass = true
					break
				}
			}
		}
		if foundSuitClass {
			break
		}
	}

	assert.True(t, foundSuitClass, "At least one card should have a card class")
}
