package views

import (
	"bytes"
	"fmt"
	"html/template"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/net/html"

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
--- Game Prompt ---
Welcome to Scopa Trainer! Click 'New Game' to start playing. [👆 New Game]
`
	actual := visualizeTemplate(doc)
	assert.Equal(t, normalizeWhitespace(expected), actual)
}

func TestGameInProgress_PlayerTurn(t *testing.T) {
	// Given a simplified game in progress model
	model := domain.NewUIModel()
	model.GameInProgress = true
	model.ShowNewGameButton = false
	model.GamePrompt = "Your turn."
	model.PlayerTurn = true

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
--- Game Prompt ---
Your turn.

--- Game Stats ---
Deck: 0 cards Your Captures: 0 cards AI Captures: 0 cards

--- Table Cards ---
Table Cards (1)
[👆 Asso-di-Coppe]

--- Player Hand ---
Your Hand (1)
[👆 Tre-di-Denari]
`
	actual := visualizeTemplate(doc)
	assert.Equal(t, normalizeWhitespace(expected), actual)
}

func TestCardSelection_WithCapturableCard(t *testing.T) {
	// Arrange
	model := domain.NewUIModel()
	model.GameInProgress = true
	model.ShowNewGameButton = false
	model.GamePrompt = "Your turn."
	model.PlayerTurn = true

	// With matching cards on table and in hand (same rank)
	model.TableCards = []domain.Card{
		{Suit: domain.Coppe, Rank: domain.Tre},
		{Suit: domain.Bastoni, Rank: domain.Quattro},
	}
	model.PlayerHand = []domain.Card{
		{Suit: domain.Denari, Rank: domain.Tre}, // This matches the Tre on the table
	}

	// Selected card that matches a table card
	model.SelectedCard = domain.Card{Suit: domain.Denari, Rank: domain.Tre}
	model.CanPlaySelectedCard = false // Should show capture message

	// Act
	doc := renderTemplate(t, model)

	// Assert
	expected := `
--- Game Prompt ---
Your turn.

--- Game Stats ---
Deck: 0 cards Your Captures: 0 cards AI Captures: 0 cards

--- Table Cards ---
Table Cards (2)
[👆 Tre-di-Coppe ⭐] [👆 Quattro-di-Bastoni] [👆 You must capture a card with the same rank]

--- Player Hand ---
Your Hand (1)
[👆 Tre-di-Denari ✓]
`
	actual := visualizeTemplate(doc)
	assert.Equal(t, normalizeWhitespace(expected), actual)
}

func TestCardSelection_CanPlayToTable(t *testing.T) {
	// Arrange
	model := domain.NewUIModel()
	model.GameInProgress = true
	model.ShowNewGameButton = false
	model.GamePrompt = "Your turn."
	model.PlayerTurn = true

	// With non-matching cards on table
	model.TableCards = []domain.Card{
		{Suit: domain.Coppe, Rank: domain.Quattro},
	}
	model.PlayerHand = []domain.Card{
		{Suit: domain.Denari, Rank: domain.Tre},
	}

	// Selected card that doesn't match any table card
	model.SelectedCard = domain.Card{Suit: domain.Denari, Rank: domain.Tre}
	model.CanPlaySelectedCard = true // Should show play to table message

	// Act
	doc := renderTemplate(t, model)

	// Assert
	expected := `
--- Game Prompt ---
Your turn.

--- Game Stats ---
Deck: 0 cards Your Captures: 0 cards AI Captures: 0 cards

--- Table Cards ---
Table Cards (1)
[👆 Quattro-di-Coppe] [👆 Click here to play the selected card to the table]

--- Player Hand ---
Your Hand (1)
[👆 Tre-di-Denari ✓]
`
	actual := visualizeTemplate(doc)
	assert.Equal(t, normalizeWhitespace(expected), actual)
}

func TestParseTemplates(t *testing.T) {
	// Test that the template parser works correctly
	templ := ParseTemplates("../../../../templates/game.html")
	require.NotNil(t, templ)

	// Verify the template has the expected name
	require.NotNil(t, templ.Lookup("game.html"), "Template should have the expected name")

	// Create a simple test template to verify the custom function
	testTemplate := `<div class="card {{.Suit | lower}}">Test</div>`
	tmpl, err := template.New("test").Funcs(template.FuncMap{
		"lower": strings.ToLower,
	}).Parse(testTemplate)
	require.NoError(t, err)

	// Test with a simple data structure
	data := struct {
		Suit string
	}{
		Suit: "COPPE",
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	require.NoError(t, err)

	// Check that the output contains the lowercase suit class
	assert.Contains(t, buf.String(), `class="card coppe"`)
}

func TestGameStates_TableDriven(t *testing.T) {
	tests := []struct {
		name     string
		model    domain.UIModel
		expected string
	}{
		{
			name: "empty table with multiple cards in hand",
			model: domain.UIModel{
				GamePrompt: "Your turn.",
				TableCards: []domain.Card{},
				PlayerHand: []domain.Card{
					{Suit: domain.Denari, Rank: domain.Re},
					{Suit: domain.Coppe, Rank: domain.Cavallo},
					{Suit: domain.Bastoni, Rank: domain.Fante},
				},
				GameInProgress:    true,
				PlayerTurn:        true,
				SelectedCard:      domain.NO_CARD_SELECTED,
				DeckCount:         30,
				ShowNewGameButton: false,
			},
			expected: `
--- Game Prompt ---
Your turn.

--- Game Stats ---
Deck: 30 cards Your Captures: 0 cards AI Captures: 0 cards

--- Table Cards ---
Table Cards (0)

--- Player Hand ---
Your Hand (3)
[👆 Re-di-Denari] [👆 Cavallo-di-Coppe] [👆 Fante-di-Bastoni]
`,
		},
		{
			name: "multiple cards on table with same rank as hand card",
			model: domain.UIModel{
				GamePrompt: "Your turn.",
				TableCards: []domain.Card{
					{Suit: domain.Denari, Rank: domain.Sette},
					{Suit: domain.Coppe, Rank: domain.Sette},
					{Suit: domain.Bastoni, Rank: domain.Due},
				},
				PlayerHand: []domain.Card{
					{Suit: domain.Spade, Rank: domain.Sette},
				},
				GameInProgress:      true,
				PlayerTurn:          true,
				SelectedCard:        domain.Card{Suit: domain.Spade, Rank: domain.Sette},
				CanPlaySelectedCard: false,
				DeckCount:           20,
				PlayerCaptureCount:  5,
				AICaptureCount:      3,
				ShowNewGameButton:   false,
			},
			expected: `
--- Game Prompt ---
Your turn.

--- Game Stats ---
Deck: 20 cards Your Captures: 5 cards AI Captures: 3 cards

--- Table Cards ---
Table Cards (3)
[👆 Sette-di-Denari ⭐] [👆 Sette-di-Coppe ⭐] [👆 Due-di-Bastoni] [👆 You must capture a card with the same rank]

--- Player Hand ---
Your Hand (1)
[👆 Sette-di-Spade ✓]
`,
		},
		{
			name: "game with captures and scores",
			model: domain.UIModel{
				GamePrompt: "AI played Fante di Spade and captured Fante di Denari.",
				TableCards: []domain.Card{
					{Suit: domain.Coppe, Rank: domain.Cinque},
					{Suit: domain.Bastoni, Rank: domain.Asso},
				},
				PlayerHand: []domain.Card{
					{Suit: domain.Spade, Rank: domain.Tre},
					{Suit: domain.Denari, Rank: domain.Cinque},
				},
				GameInProgress:     true,
				PlayerTurn:         true,
				SelectedCard:       domain.NO_CARD_SELECTED,
				DeckCount:          10,
				PlayerCaptureCount: 12,
				AICaptureCount:     14,
				ShowNewGameButton:  false,
			},
			expected: `
--- Game Prompt ---
AI played Fante di Spade and captured Fante di Denari.

--- Game Stats ---
Deck: 10 cards Your Captures: 12 cards AI Captures: 14 cards

--- Table Cards ---
Table Cards (2)
[👆 Cinque-di-Coppe] [👆 Asso-di-Bastoni]

--- Player Hand ---
Your Hand (2)
[👆 Tre-di-Spade] [👆 Cinque-di-Denari]
`,
		},
		{
			name: "selected card with capturable cards",
			model: domain.UIModel{
				GamePrompt: "Your turn.",
				TableCards: []domain.Card{
					{Suit: domain.Coppe, Rank: domain.Sette},
					{Suit: domain.Denari, Rank: domain.Sette},
					{Suit: domain.Bastoni, Rank: domain.Due},
				},
				PlayerHand: []domain.Card{
					{Suit: domain.Spade, Rank: domain.Sette},
					{Suit: domain.Denari, Rank: domain.Re},
				},
				GameInProgress:      true,
				PlayerTurn:          true,
				SelectedCard:        domain.Card{Suit: domain.Spade, Rank: domain.Sette},
				CanPlaySelectedCard: false,
				DeckCount:           15,
				PlayerCaptureCount:  8,
				AICaptureCount:      7,
				ShowNewGameButton:   false,
			},
			expected: `
--- Game Prompt ---
Your turn.

--- Game Stats ---
Deck: 15 cards Your Captures: 8 cards AI Captures: 7 cards

--- Table Cards ---
Table Cards (3)
[👆 Sette-di-Coppe ⭐] [👆 Sette-di-Denari ⭐] [👆 Due-di-Bastoni] [👆 You must capture a card with the same rank]

--- Player Hand ---
Your Hand (2)
[👆 Sette-di-Spade ✓] [👆 Re-di-Denari]
`,
		},
		{
			name: "AI turn with disabled player buttons",
			model: domain.UIModel{
				GamePrompt: "AI is thinking...",
				TableCards: []domain.Card{
					{Suit: domain.Coppe, Rank: domain.Sette},
					{Suit: domain.Denari, Rank: domain.Tre},
				},
				PlayerHand: []domain.Card{
					{Suit: domain.Spade, Rank: domain.Quattro},
					{Suit: domain.Bastoni, Rank: domain.Cinque},
				},
				GameInProgress:    true,
				PlayerTurn:        false,
				SelectedCard:      domain.NO_CARD_SELECTED,
				ShowNewGameButton: false,
			},
			expected: `
--- Game Prompt ---
AI is thinking...

--- Game Stats ---
Deck: 0 cards Your Captures: 0 cards AI Captures: 0 cards

--- AI Turn ---
[👆 Let AI Play Its Turn]

--- Table Cards ---
Table Cards (2)
Sette-di-Coppe disabled Tre-di-Denari disabled

--- Player Hand ---
Your Hand (2)
Quattro-di-Spade disabled Cinque-di-Bastoni disabled
`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Act
			doc := renderTemplate(t, test.model)

			// Assert
			actual := visualizeTemplate(doc)
			assert.Equal(t, normalizeWhitespace(test.expected), actual)
		})
	}
}

func TestDisabledPlayArea(t *testing.T) {
	// Arrange - Test when play area is disabled
	model := domain.UIModel{
		GamePrompt: "Your turn.",
		TableCards: []domain.Card{
			{Suit: domain.Coppe, Rank: domain.Sette},
		},
		PlayerHand: []domain.Card{
			{Suit: domain.Spade, Rank: domain.Tre},
		},
		GameInProgress:      true,
		PlayerTurn:          true,
		SelectedCard:        domain.Card{Suit: domain.Spade, Rank: domain.Tre},
		CanPlaySelectedCard: false,
		ShowNewGameButton:   false,
	}

	// Act
	doc := renderTemplate(t, model)

	// Assert
	assert.Contains(t, doc, `class="play-area disabled"`)
	assert.Contains(t, doc, `You must capture a card with the same rank`)
}

func renderTemplate(t *testing.T, model domain.UIModel) string {
	templ := ParseTemplates("../../../../templates/game.html")

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
	// First, replace all whitespace sequences within lines with a single space
	inlineNormalized := replaceAll(s, `[^\S\n]+`, " ")

	// Then, replace multiple newlines with a single newline
	multiNewlineNormalized := replaceAll(inlineNormalized, `\n+`, "\n")

	// Finally, trim leading/trailing whitespace from each line
	lines := strings.Split(multiNewlineNormalized, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimSpace(line)
	}

	return strings.TrimSpace(strings.Join(lines, "\n"))
}

// visualizeTemplate returns a string representation of the HTML content
// that is easier to read and test.
// see https://martinfowler.com/articles/tdd-html-templates.html
// We add a `data-test-icon` attribute to some elements to simplify their rendering.
// This is a workaround to avoid having to deal with the complexity of the actual HTML structure,
// but still lets us test the logic in the template.
// This trick I learned from Esko Luontola.
func visualizeTemplate(htmlContent string) string {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return fmt.Sprintf("Error parsing HTML: %v", err)
	}

	var output strings.Builder
	skipElements := []string{"style", "script", "head"}
	visualizeNode(doc, &output, 0, skipElements)

	return normalizeWhitespace(output.String())
}

func visualizeNode(n *html.Node, output *strings.Builder, depth int, skipElements []string) {
	// Skip elements in the skip list
	if n.Type == html.ElementNode && slices.Contains(skipElements, n.Data) {
		return
	}

	// Add newlines for major UI sections
	if n.Type == html.ElementNode && n.Data == "section" {
		// Get section name from data-section attribute
		var sectionName string
		for _, attr := range n.Attr {
			if attr.Key == "data-section" {
				sectionName = attr.Val
				break
			}
		}

		// Add a newline before the section
		output.WriteString("\n")

		// If section has a name, add it as a header
		if sectionName != "" {
			output.WriteString(fmt.Sprintf("--- %s ---\n", sectionName))
		}

		// Process children
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			visualizeNode(c, output, depth+1, skipElements)
		}

		return
	}

	if n.Type == html.ElementNode && n.Data == "button" {
		output.WriteString(fmt.Sprintf("[👆 %s]\n", extractTextContent(n)))
		return
	}

	// Special handling for play area
	if n.Type == html.ElementNode && n.Data == "div" {
		var isPlayArea bool
		for _, attr := range n.Attr {
			if attr.Key == "class" && strings.Contains(attr.Val, "play-area") {
				isPlayArea = true
				break
			}
		}

		if isPlayArea {
			// Extract the text content
			text := extractTextContent(n)
			output.WriteString(fmt.Sprintf("[👆 %s]\n", text))
			return
		}
	}

	if n.Type == html.ElementNode {
		var testIcon string
		var hasOnClick bool
		var onClickValue string
		var isSelected bool
		var isCapturable bool
		var isDisabled bool

		for _, attr := range n.Attr {
			if attr.Key == "data-test-icon" {
				testIcon = attr.Val
			}
			if attr.Key == "onclick" {
				hasOnClick = true
				onClickValue = attr.Val
				_ = onClickValue
			}
			if attr.Key == "class" {
				if strings.Contains(attr.Val, "selected") {
					isSelected = true
				}
				if strings.Contains(attr.Val, "capturable") {
					isCapturable = true
				}
				if strings.Contains(attr.Val, "disabled") {
					isDisabled = true
				}
			}
		}

		if testIcon != "" && hasOnClick {
			// Add visual indicators for selected and capturable states
			var indicators string
			if isSelected {
				indicators += "✓" // Checkmark for selected
			}
			if isCapturable {
				indicators += "⭐" // Star for capturable
			}
			if isDisabled {
				indicators += " disabled" // Text indicator for disabled
			}

			if indicators != "" {
				output.WriteString(fmt.Sprintf("[👆 %s %s] ", testIcon, indicators))
			} else {
				output.WriteString(fmt.Sprintf("[👆 %s] ", testIcon))
			}
			return
		} else if testIcon != "" {
			if isDisabled {
				output.WriteString(fmt.Sprintf("%s disabled ", testIcon))
			} else {
				output.WriteString(testIcon)
				output.WriteString(" ")
			}
			return
		} else if hasOnClick {
			textContent := extractTextContent(n)
			if isDisabled {
				output.WriteString(fmt.Sprintf("[👆 %s disabled] ", textContent))
			} else {
				output.WriteString(fmt.Sprintf("[👆 %s] ", textContent))
			}
			return
		}
	} else if n.Type == html.TextNode {
		text := strings.TrimSpace(n.Data)
		if text != "" {
			// Check if this is a heading (h1, h2, etc.)
			if n.Parent != nil && n.Parent.Type == html.ElementNode &&
				(n.Parent.Data == "h1" || n.Parent.Data == "h2" || n.Parent.Data == "h3") {
				output.WriteString(text + "\n")
			} else {
				output.WriteString(text + " ")
			}
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
