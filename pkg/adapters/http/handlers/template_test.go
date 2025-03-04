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
	// Given a game in progress model
	model := domain.NewUIModel()
	model.GameInProgress = true
	model.ShowNewGameButton = false
	model.GamePrompt = "Your turn. Select a card to play."

	// With cards on the table
	model.TableCards = []domain.Card{
		{Suit: domain.Coppe, Rank: domain.Asso, Name: "Asso", Value: 1},
	}

	// And cards in the player's hand
	model.PlayerHand = []domain.Card{
		{Suit: domain.Denari, Rank: domain.Tre, Name: "Tre", Value: 3},
		{Suit: domain.Bastoni, Rank: domain.Re, Name: "Re", Value: 10},
	}

	// When rendering the template
	doc := renderTemplate(t, model)

	// Then we should see the game prompt
	promptDiv := findElement(doc, "div", "game-prompt")
	assert.NotNil(t, promptDiv)
	assert.Contains(t, getTextContent(promptDiv), "Your turn")

	// And we should not see the new game button
	assert.Nil(t, findElement(doc, "button", "new-game-button"))

	// And we should see the table area with one card
	tableArea := findElement(doc, "div", "table-area")
	assert.NotNil(t, tableArea)
	tableHeading := findElement(tableArea, "h2", "")
	assert.Contains(t, getTextContent(tableHeading), "Table Cards (1)")

	tableCards := findAllElements(tableArea, "div", "card")
	assert.Len(t, tableCards, 1)
	assert.Contains(t, getTextContent(tableCards[0]), "Asso di Coppe")

	// And we should see the player's hand with two cards
	handArea := findElement(doc, "div", "hand-area")
	assert.NotNil(t, handArea)
	handHeading := findElement(handArea, "h2", "")
	assert.Contains(t, getTextContent(handHeading), "Your Hand (2)")

	handCards := findAllElements(handArea, "div", "card")
	assert.Len(t, handCards, 2)

	// Find card contents (order doesn't matter)
	cardTexts := []string{
		getTextContent(handCards[0]),
		getTextContent(handCards[1]),
	}
	assert.Contains(t, cardTexts, "Tre di Denari")
	assert.Contains(t, cardTexts, "Re di Bastoni")
}

func TestCardSuitStyling(t *testing.T) {
	// Given a model with cards of each suit
	model := domain.NewUIModel()
	model.GameInProgress = true
	model.TableCards = []domain.Card{
		{Suit: domain.Coppe, Rank: domain.Asso, Name: "Asso", Value: 1},
		{Suit: domain.Denari, Rank: domain.Due, Name: "Due", Value: 2},
		{Suit: domain.Bastoni, Rank: domain.Tre, Name: "Tre", Value: 3},
		{Suit: domain.Spade, Rank: domain.Quattro, Name: "Quattro", Value: 4},
	}

	// When rendering the template
	doc := renderTemplate(t, model)

	// Then each card should have the appropriate suit class
	cards := findAllElements(doc, "div", "card")
	assert.Len(t, cards, 14)

	// Verify one card of each suit class exists
	suitClasses := map[string]bool{
		"coppe":   false,
		"denari":  false,
		"bastoni": false,
		"spade":   false,
	}

	for _, card := range cards {
		for _, attr := range card.Attr {
			if attr.Key == "class" {
				for suit := range suitClasses {
					if strings.Contains(attr.Val, suit) {
						suitClasses[suit] = true
					}
				}
			}
		}
	}

	for suit, found := range suitClasses {
		assert.True(t, found, "Should have a card with %s class", suit)
	}
}
