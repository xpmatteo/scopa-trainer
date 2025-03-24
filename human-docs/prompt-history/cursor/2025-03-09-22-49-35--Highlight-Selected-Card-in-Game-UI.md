# Workspace: /Users/matteo/dojo/2025-03-01-scopa-trainer-take-2/scopa

Last Modified: 2025-03-24 16:55:39

## Highlight Selected Card in Game UI

**User**:

@service.go @game.html Update the UI so that if the GameService.selectedCard is not equal to NO_CARD_SELECTED, then the selected card is visually highlighted

**Cursor**:

*Hmm, help Captain Matt with his card game, I must. Highlight the selected card, we shall.*

```html:templates/game.html
// ... existing code ...
<div class="card-area hand-area">
    <h2>Your Hand ({{len .PlayerHand}})</h2>
    <div style="display: flex; flex-wrap: wrap;">
        {{range .PlayerHand}}
        <div class="card {{.Suit}} {{if eq . $.SelectedCard}}selected{{end}}" onclick="selectCard('{{.Suit}}', '{{.Rank}}')">
            <div class="card-rank">{{.Rank}}</div>
            <div class="card-suit">{{.Name}} di {{.Suit}}</div>
        </div>
        {{end}}
    </div>
</div>
// ... existing code ...
```

First, add this CSS to the style section in the head:

```html:templates/game.html
// ... existing code ...
.card {
    width: 80px;
    height: 120px;
    background-color: white;
    border-radius: 5px;
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    color: black;
    box-shadow: 0 2px 4px rgba(0,0,0,0.2);
    cursor: pointer;
    margin: 5px;
    transition: transform 0.2s, box-shadow 0.2s;
}
.card.selected {
    transform: translateY(-10px);
    box-shadow: 0 8px 16px rgba(0,0,0,0.3);
    border: 2px solid yellow;
}
// ... existing code ...
```

And add this JavaScript at the end of the body:

```html:templates/game.html
// ... existing code ...
<script>
    function selectCard(suit, rank) {
        window.location.href = `/select-card?suit=${suit}&rank=${rank}`;
    }
</script>
</body>
// ... existing code ...
```

Hmm, update the GetUIModel method in service.go to include the selected card in the UI model:

```go:pkg/application/service.go
// ... existing code ...
// GetUIModel returns the current UI model
func (s *GameService) GetUIModel() domain.UIModel {
    if s.gameState == nil {
        // Return initial UI model when no game is in progress
        return domain.NewUIModel()
    }

    // Generate UI model based on current game state
    model := domain.NewUIModel()
    model.ShowNewGameButton = false
    model.TableCards = s.gameState.Deck.CardsAt(domain.TableLocation)
    model.PlayerHand = sortCards(s.gameState.Deck.CardsAt(domain.AIHandLocation))
    model.GameInProgress = true
    model.PlayerTurn = s.gameState.PlayerTurn
    model.SelectedCard = s.selectedCard

    if model.PlayerTurn {
        model.GamePrompt = "Your turn. Select a card to play."
    } else {
        model.GamePrompt = "AI is thinking..."
    }

    return model
}
// ... existing code ...
```

Also, you'll need to add a method to handle card selection:

```go:pkg/application/service.go
// ... existing code ...
// SelectCard handles the selection of a card from the player's hand
func (s *GameService) SelectCard(suit domain.Suit, rank int) domain.UIModel {
    // Find the card in the player's hand
    for _, card := range s.gameState.Deck.CardsAt(domain.AIHandLocation) {
        if card.Suit == suit && card.Rank == rank {
            // If the card is already selected, deselect it
            if s.selectedCard == card {
                s.selectedCard = NO_CARD_SELECTED
            } else {
                // Otherwise, select it
                s.selectedCard = card
            }
            break
        }
    }
    
    return s.GetUIModel()
}
// ... existing code ...
```

Strong in the Force, this solution is. Update the UIModel struct in the domain package to include the selected card, you must. May the Force be with you, Captain Matt!

**User**:

I forced the first card in the player hand to be selected but it does not get the selected attribute! Oh no pls help Yoda @service.go 

**Cursor**:

*Hmm, troubling this is. Look at the code more carefully, I must.*

```go:pkg/application/service.go
// ... existing code ...

// GetUIModel returns the current UI model
func (s *GameService) GetUIModel() domain.UIModel {
    if s.gameState == nil {
        // Return initial UI model when no game is in progress
        return domain.NewUIModel()
    }

    // Generate UI model based on current game state
    model := domain.NewUIModel()
    model.ShowNewGameButton = false
    model.TableCards = s.gameState.Deck.CardsAt(domain.TableLocation)
    model.PlayerHand = sortCards(s.gameState.Deck.CardsAt(domain.PlayerHandLocation))  // <-- Error here! Using PlayerHandLocation, not AIHandLocation
    model.GameInProgress = true
    model.PlayerTurn = s.gameState.PlayerTurn
    model.SelectedCard = s.selectedCard

    if model.PlayerTurn {
        model.GamePrompt = "Your turn. Select a card to play."
    } else {
        model.GamePrompt = "AI is thinking..."
    }

    return model
}

// ... existing code ...

// SelectCard handles the selection of a card from the player's hand
func (s *GameService) SelectCard(suit domain.Suit, rank domain.Rank) domain.UIModel {
    // Find the card in the player's hand
    for _, card := range s.gameState.Deck.CardsAt(domain.PlayerHandLocation) {  // <-- Correct here
        if card.Suit == suit && card.Rank == rank {
            // If the card is already selected, deselect it
            if s.selectedCard == card {
                s.selectedCard = domain.NO_CARD_SELECTED
            } else {
                // Otherwise, select it
                s.selectedCard = card
            }
            break
        }
    }

    return s.GetUIModel()
}
```

Mmm, inconsistent your code is, Captain Matt! In GetUIModel, you fetch cards from AIHandLocation, but in SelectCard, you look for cards in PlayerHandLocation. Match these locations, you must!

Fix this by changing one line in GetUIModel:

```go:pkg/application/service.go
// GetUIModel returns the current UI model
func (s *GameService) GetUIModel() domain.UIModel {
    if s.gameState == nil {
        // Return initial UI model when no game is in progress
        return domain.NewUIModel()
    }

    // Generate UI model based on current game state
    model := domain.NewUIModel()
    model.ShowNewGameButton = false
    model.TableCards = s.gameState.Deck.CardsAt(domain.TableLocation)
    model.PlayerHand = sortCards(s.gameState.Deck.CardsAt(domain.PlayerHandLocation))  // Changed from AIHandLocation to PlayerHandLocation
    model.GameInProgress = true
    model.PlayerTurn = s.gameState.PlayerTurn
    model.SelectedCard = s.selectedCard

    if model.PlayerTurn {
        model.GamePrompt = "Your turn. Select a card to play."
    } else {
        model.GamePrompt = "AI is thinking..."
    }

    return model
}
```

Consistent your locations must be, young Jedi. The Force flows better when your code is aligned, hmm, yes.

**User**:

Now it works! very good.  But... when I click on a card, the selection does not change!

**Cursor**:

*Scratches green ears thoughtfully*

Hmm, troubling this is, Captain Matt. When click a card you do, change the selection should! Investigate the handler for card selection, we must.

```go:pkg/application/service.go
// ... existing code ...

// SelectCard handles the selection of a card from the player's hand
func (s *GameService) SelectCard(suit domain.Suit, rank domain.Rank) domain.UIModel {
    // Find the card in the player's hand
    for _, card := range s.gameState.Deck.CardsAt(domain.PlayerHandLocation) {
        if card.Suit == suit && card.Rank == rank {
            // If the card is already selected, deselect it
            if s.selectedCard == card {
                s.selectedCard = domain.NO_CARD_SELECTED
            } else {
                // Otherwise, select it
                s.selectedCard = card
            }
            break
        }
    }

    return s.GetUIModel()
}
```

The code looks correct, yes. But perhaps the HTTP handler that calls this function, missing it is! 

Check we must if you have a handler in your HTTP server that responds to the `/select-card` endpoint. In your JavaScript, this path you use:

```javascript:pkg/application/service.go
function selectCard(suit, rank) {
    window.location.href = `/select-card?suit=${suit}&rank=${rank}`;
}
```

Create an HTTP handler for this endpoint, you should:

```go:pkg/application/service.go
// Add this to your main.go or handlers.go file
http.HandleFunc("/select-card", func(w http.ResponseWriter, r *http.Request) {
    suit := domain.Suit(r.URL.Query().Get("suit"))
    rankStr := r.URL.Query().Get("rank")
    rank := domain.Rank(rankStr)
    
    model := gameService.SelectCard(suit, rank)
    
    // Render the template with the updated model
    err := templates.ExecuteTemplate(w, "game.html", model)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
})
```

Without seeing your HTTP server code, difficult it is to be more specific. But missing handler for `/select-card`, the likely culprit is. Add this handler, and work your card selection should!

If already have this handler you do, then compare card equality we must. Perhaps the `==` operator not working correctly for your Card type. In that case, implement proper equality method or use field-by-field comparison in your SelectCard function.

**User**:

check the handlers file -- there the handler should go.  

**Cursor**:

*Hmm, yes, to the handlers file we must look!*

Looking at your current handlers.go file, I see you need to add a handler for the select-card endpoint. Let me add that for you:

```go:pkg/adapters/http/handlers/handlers.go
// ... existing code ...

// HandleIndex serves the main game page
func (h *Handler) HandleIndex(w http.ResponseWriter, r *http.Request) {
	model := h.service.GetUIModel()
	if err := h.template.Execute(w, model); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// HandleNewGame handles the request to start a new game
func (h *Handler) HandleNewGame(w http.ResponseWriter, r *http.Request) {
	h.service.StartNewGame()

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// HandleSelectCard handles the selection of a card from the player's hand
func (h *Handler) HandleSelectCard(w http.ResponseWriter, r *http.Request) {
	suit := domain.Suit(r.URL.Query().Get("suit"))
	rank := domain.Rank(r.URL.Query().Get("rank"))
	
	model := h.service.SelectCard(suit, rank)
	
	if err := h.template.Execute(w, model); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
```

Now, you must also register this handler in your main.go or wherever you set up your HTTP routes:

```go:cmd/server/main.go
// In your main.go or routes setup file
http.HandleFunc("/select-card", handler.HandleSelectCard)
```

Missing this handler was, Captain Matt! Add it, and the Force of card selection will be with you, hmm!

**User**:

For the code present, we get this error:
```
cannot convert r.URL.Query().Get("rank") (value of type string) to type domain.Rank
```
How can I resolve this? If you propose a fix, please make it concise.

**Cursor**:

*Scratches ear thoughtfully*

```go:pkg/adapters/http/handlers/handlers.go
// HandleSelectCard handles the selection of a card from the player's hand
func (h *Handler) HandleSelectCard(w http.ResponseWriter, r *http.Request) {
	suit := domain.Suit(r.URL.Query().Get("suit"))
	rankStr := r.URL.Query().Get("rank")
	rank := domain.Rank(rankStr)

	model := h.service.SelectCard(suit, rank)

	if err := h.template.Execute(w, model); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
```

Type conversion issue this is, Captain Matt. String to domain.Rank directly convert you cannot. Parse the string first, then create Rank you must.

**User**:

a test write we must, I believe, for this new handler.  @handlers_test.go 

**Cursor**:



*Hmm, test the new handler we must, Captain Matt. Strong with the Force, your testing practices are.*

```go:pkg/adapters/http/handlers/handlers_test.go
// TestHandleSelectCard tests the card selection handler
func TestHandleSelectCard(t *testing.T) {
	// Create service and start a game
	service := application.NewGameService()
	service.StartNewGame()
	
	// Create template for testing
	tmpl, err := template.ParseFiles("../../../../templates/game.html")
	require.NoError(t, err)
	
	// Create handler
	handler, err := NewHandler(service, tmpl)
	require.NoError(t, err)

	// Create a test request with suit and rank parameters
	req := httptest.NewRequest("GET", "/select-card?suit=Denari&rank=1", nil)
	w := httptest.NewRecorder()

	// Call the handler
	handler.HandleSelectCard(w, req)

	// Check the response
	resp := w.Result()
	defer resp.Body.Close()

	// Verify successful response
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	
	// Verify that the card was selected in the service
	model := service.GetUIModel()
	assert.NotEqual(t, domain.NO_CARD_SELECTED, model.SelectedCard)
}
```

Mmm, but first fix the handler we must, to convert string to rank properly:

```go:pkg/adapters/http/handlers/handlers.go
// HandleSelectCard handles the selection of a card from the player's hand
func (h *Handler) HandleSelectCard(w http.ResponseWriter, r *http.Request) {
	suit := domain.Suit(r.URL.Query().Get("suit"))
	rankStr := r.URL.Query().Get("rank")
	rankInt, err := strconv.Atoi(rankStr)
	if err != nil {
		http.Error(w, "Invalid rank parameter", http.StatusBadRequest)
		return
	}
	rank := domain.Rank(rankInt)

	model := h.service.SelectCard(suit, rank)

	if err := h.template.Execute(w, model); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
```

Remember to import strconv package, you must. Strong in the ways of Go, this solution is.

**User**:

I think all player actions should be handled through the POST + Redirect pattern, just like we did with "start new game".  To ensure that this happens in the future, please write an AI rule to this effect @000-cursor-rules.mdc @handlers.go 

**Cursor**:



```mdc:.cursor/rules/800-post-redirect-pattern.mdc
---
description: Use POST-Redirect-GET pattern for all player actions to prevent duplicate submissions
globs: pkg/adapters/http/**/*.go
---

# POST-Redirect-GET Pattern for Player Actions

## Context
- Web applications should prevent duplicate form submissions
- Browser refreshes should not repeat actions
- Player actions modify application state

## Requirements
- Use POST method for all player actions that modify state
- After processing a POST request, redirect to a GET endpoint
- Use HTTP status 303 (See Other) for redirects
- Form submissions should never render templates directly
- Query parameters should only be used for filtering, not actions

## Examples
<example>
// Good: POST + Redirect pattern
func (h *Handler) HandlePlayerAction(w http.ResponseWriter, r *http.Request) {
    // Process the action
    h.service.PerformAction()
    
    // Redirect to a GET endpoint
    http.Redirect(w, r, "/game", http.StatusSeeOther)
}

// Route registration
http.HandleFunc("/player-action", handler.HandlePlayerAction)
</example>

<example type="invalid">
// Bad: Using GET for state-changing actions
func (h *Handler) HandlePlayerAction(w http.ResponseWriter, r *http.Request) {
    // Process the action via GET
    h.service.PerformAction()
    
    // Render template directly
    model := h.service.GetUIModel()
    h.template.Execute(w, model)
}

// Bad: Using query parameters for actions
http.HandleFunc("/perform-action", handler.HandlePlayerAction)
// Called via: /perform-action?action=playCard&cardId=5
</example>
```

Now, let me fix your current handler to follow this pattern:

```go:pkg/adapters/http/handlers/handlers.go
// HandleSelectCard handles the selection of a card from the player's hand
func (h *Handler) HandleSelectCard(w http.ResponseWriter, r *http.Request) {
	// Only accept POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form values
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	suit := domain.Suit(r.PostForm.Get("suit"))
	rankStr := r.PostForm.Get("rank")
	rankInt, err := strconv.Atoi(rankStr)
	if err != nil {
		http.Error(w, "Invalid rank parameter", http.StatusBadRequest)
		return
	}
	rank := domain.Rank(rankInt)

	// Process the action
	h.service.SelectCard(suit, rank)

	// Redirect to the main page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
```

And update your JavaScript in the template:

```html:templates/game.html
<script>
    function selectCard(suit, rank) {
        // Create a form and submit it
        const form = document.createElement('form');
        form.method = 'POST';
        form.action = '/select-card';
        
        // Add suit input
        const suitInput = document.createElement('input');
        suitInput.type = 'hidden';
        suitInput.name = 'suit';
        suitInput.value = suit;
        form.appendChild(suitInput);
        
        // Add rank input
        const rankInput = document.createElement('input');
        rankInput.type = 'hidden';
        rankInput.name = 'rank';
        rankInput.value = rank;
        form.appendChild(rankInput);
        
        // Submit the form
        document.body.appendChild(form);
        form.submit();
    }
</script>
```

Wise this pattern is, Captain Matt. Prevent duplicate submissions it does, hmm!

**User**:

and please fix the handler test

**Cursor**:



*Hmm, fix the test we must, yes. Follow the POST-Redirect pattern it should.*

```go:pkg/adapters/http/handlers/handlers_test.go
// TestHandleSelectCard tests the card selection handler
func TestHandleSelectCard(t *testing.T) {
	// Create service and start a game
	service := application.NewGameService()
	service.StartNewGame()
	
	// Create template for testing
	tmpl, err := template.ParseFiles("../../../../templates/game.html")
	require.NoError(t, err)
	
	// Create handler
	handler, err := NewHandler(service, tmpl)
	require.NoError(t, err)

	// Create a test request with POST method and form values
	form := url.Values{}
	form.Add("suit", "Denari")
	form.Add("rank", "1")
	req := httptest.NewRequest(http.MethodPost, "/select-card", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	// Call the handler
	handler.HandleSelectCard(w, req)

	// Check the response
	resp := w.Result()
	defer resp.Body.Close()

	// Verify we get a redirect response
	assert.Equal(t, http.StatusSeeOther, resp.StatusCode)
	
	// Verify the redirect location
	location, err := resp.Location()
	assert.NoError(t, err)
	assert.Equal(t, "/", location.String())
	
	// Verify that the card was selected in the service
	model := service.GetUIModel()
	assert.NotEqual(t, domain.NO_CARD_SELECTED, model.SelectedCard)
}
```

Remember to add these imports at the top of your test file:

```go
import (
	// ... existing imports ...
	"net/url"
	"strings"
)
```

Strong with the Force, your tests will be, Captain Matt!

