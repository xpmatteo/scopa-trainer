# Scopa Trainer User Stories

## Epics

### Epic 1: Game Setup
When a player opens the application, they can set up and start a new Scopa game to play against the AI.

### Epic 2: Core Gameplay
When a player engages with the game, they can make moves according to Scopa rules and play a complete game against the AI.

### Epic 3: Game State Management
When a player interacts with the game, their progress is automatically saved and can be resumed later.

### Epic 4: Post-Game Analysis
When a player completes a game, they receive analysis and feedback to improve their skills.

## User Stories

### Game Setup

✅ **US1.1: Start New Game**
When a player opens the application, they can start a new game of Scopa to begin playing against the AI opponent.
Acceptance Criteria:
1. A "New Game" button is displayed on the main screen
2. Clicking the button initializes a fresh game with a standard 40-card Italian deck
3. The game state shows 3-4 cards on the table and 3 cards in the player's hand
4. Any previous game state is cleared from storage

**US1.2: Choose First Player**
When a player starts a new game, they can choose whether they or the AI plays first.
Acceptance Criteria:
1. Before the game begins, an option to select who goes first is presented
2. The selection affects who takes the first turn in the game
3. The turn indicator correctly shows the selected first player

### Core Gameplay

✅ **US2.1: View Game State** 
When a player looks at the game screen, they can see all relevant game information.
Acceptance Criteria:
1. Player's hand is visible at the bottom of the screen
2. Table cards are displayed in the center
3. AI's hand appears as face-down cards at the top
4. Current score, remaining deck count, and turn indicator are visible

✅ **US2.2: Select Card from Hand**
When it's the player's turn, they can select a card from their hand to play.
Acceptance Criteria:
1. Player can click on any card in their hand to select it
2. Selected card is visually highlighted
3. Playing on a different card selects it, and deselects previous one

**US2.3: Capture Single Card**
When a player selects a card from their hand, they can capture a matching card from the table.
Acceptance Criteria:
1. After selecting a hand card, player can click on a table card of matching value
2. The capture removes both cards from display and adds them to player's capture pile

**US2.4: Capture Card Combinations**
When a player selects a card from their hand, they can capture multiple cards that sum to its value.
Acceptance Criteria:
1. Player can select multiple table cards that sum to their played card's value
2. Selected combination is highlighted before confirming
3. Invalid combinations cannot be selected

**US2.5: Select from Multiple Capture Options**
When multiple valid capture combinations exist, the player can choose which one to take.
Acceptance Criteria:
1. All valid capture combinations are presented to the player
2. Player can select which combination to capture
3. Only one combination can be selected at a time

**US2.6: Play Card to Table**
When a player cannot or chooses not to make a capture, they can play a card to the table.
Acceptance Criteria:
1. If player selects a hand card and no capture is possible, the selected card is added to the table cards
2. Turn passes to the AI after card is played

**US2.7: Experience AI's Turn**
When it's the AI's turn, the player can observe the AI opponent making its move.
Acceptance Criteria:
1. AI's move is visually indicated with appropriate animations
2. Captures made by the AI are clearly shown
3. AI's turn completes automatically within 2 seconds

**US2.8: Receive New Cards**
When both players have played all cards in their hands, they receive new cards if available.
Acceptance Criteria:
1. When both hands are empty, 3 new cards are dealt to each player
2. The dealing process is visible to the player
3. If no cards remain in the deck, the game proceeds to final scoring

**US2.9: End Game and View Score**
When the deck is exhausted and all cards have been played, the player sees the final score.
Acceptance Criteria:
1. A "Game Over" message appears when the game ends
2. Final scores for both player and AI are displayed
3. Breakdown of points (scopa, primiera) is shown
4. Options to analyze game or start a new one are presented

### Game State Management

**US3.1: Auto-Save Game Progress**
When the player makes a move, the game state is automatically saved.
Acceptance Criteria:
1. After each move, current game state is saved to localStorage
2. Saving process is invisible to the player
3. Saved state includes all information needed to resume the exact game state

**US3.2: Resume Saved Game**
When the player returns to the application, they can continue a previously saved game.
Acceptance Criteria:
1. If a saved game exists, an option to resume it is presented
2. Resumed game restores the exact previous state
3. All game elements (scores, cards, turn indicator) are correctly restored

**US3.3: Undo Move in New Game**
When the player makes the first move in a new game, they can undo it before the AI responds.
Acceptance Criteria:
1. An "Undo" button is available after the player's first move
2. Clicking it reverts the game state to before the move
3. The undo feature is disabled after the AI plays its turn

### Post-Game Analysis

**US4.1: View Game Summary**
When a game ends, the player can see a summary of the game performance.
Acceptance Criteria:
1. Summary shows final scores and key statistics
2. Summary includes capture counts and point breakdowns
3. A clear indication of who won is displayed

**US4.2: Review Move Quality**
When viewing post-game analysis, the player can see an evaluation of each of their moves.
Acceptance Criteria:
1. Each move is rated (perfect, suboptimal, terrible)
2. Ratings consider only information available to player at the time
3. A count of moves in each quality category is provided

**US4.3: Analyze Key Decision Points**
When reviewing the game, the player can focus on critical moves that affected the outcome.
Acceptance Criteria:
1. Key decision points are highlighted in the game replay
2. For each key decision, alternative moves and potential outcomes are shown
3. Text explanations of AI reasoning are provided

**US4.4: Receive Improvement Recommendations**
When the post-game analysis completes, the player receives personalized tips.
Acceptance Criteria:
1. Specific improvement recommendations based on player's moves are provided
2. Recommendations address patterns in decision-making
3. Tips are actionable and specific to Scopa strategy

**US4.5: Replay Key Decisions**
When analyzing the game, the player can interactively replay important moments.
Acceptance Criteria:
1. Player can select specific points in the game to replay
2. Game state is reconstructed as it was at that time
3. Player can try alternative moves to see outcomes
4. AI Tutor provides feedback on alternative choices