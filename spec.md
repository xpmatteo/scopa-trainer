# Scopa Trainer Specification

## Overview
A browser-based implementation of the Scopa card game (two-player variant) that allows a human player to play against an AI opponent and receive post-game analysis to improve their skills.

## Actors

- Player: the human user
- AI Opponent: plays against the Player
- AI Tutor: provides post-game analysis
- Server: provides the mechanisms 

## Game Rules

### Card Deck
- Use an Italian 40 cards deck
- Italian suits: Coppe, Denari, Bastoni, Spade
- Cards:
  - "Asso" ranks 1
  - "Fante" ranks 8
  - "Cavallo" ranks 9
  - "Re" ranks 10
  - All other cards rank as their face value

### Game Flow
1. Each player is dealt 10 cards
2. Players take turns playing one card at a time
3. When both players have played all cards, they are each dealt 10 more cards
4. Game ends when the deck is exhausted

### Capturing Mechanics
- Player selects a card from their hand, then selects a matching card or combination on the table
- Player must explicitly choose which combination to capture if multiple options exist
- If a player cannot make a capture, they just play a card to the table
- Clearing all cards from the table (scopa) earns a point

### Scoring
- Simplified primiera: 1 point for capturing 3-4 sevens; no point if players capture same number
- Scopa: 1 point for each time a player clears the table
- Final score displayed at game conclusion

## Technical Architecture

### Technology Stack
- Go programming language
- Simple HTML/CSS for UI
- Local browser storage for game state

### Code Organization
- Core game logic in Go (rules, state management, AI)
- Use standard Go project layout
- UI rendering in HTML/CSS using the standard Go http templates
- Use the hexagonal architecture: the http handlers are adapters; they call application services inside the hexagon
- Application services implement user actions and return the complete updated UI model

## User Interface

### Layout
- Standard computerized card game layout
- Player's hand at bottom of screen
- Table cards in center
- AI's hand represented as face-down cards at top
- Score display visible during gameplay
- Turn indicator at top of screen (text prompt)

### Card Display
- Cards shown with standard illustrations
- No special highlighting for potential captures
- Selected cards should have visual indication

### Game State Display
- Current score always visible
- Capture count for both players visible
- Number of remaining cards in deck visible

## AI Implementation

### Strategy
- AI keeps track of cards it has seen played
- Uses minimax algorithm to determine optimal moves
- Does not "cheat" by looking at hidden cards
- Can provide text explanations of its reasoning when requested

## Game Features

### Player Controls
- Card selection: first hand card, then table card(s)
- No computer assistance identifying potential captures
- Option to select who plays first at game start
- Game restart with undo capability until first move of new game

### State Management
- Current game saved in localStorage
- Starting new game erases previous game state
- Full game state recorded for each move to enable replay

## Post-Game Analysis

### Analysis Features
- Turn-by-turn evaluation of move quality (perfect, suboptimal, terrible)
- Assessment based on what player could know at the time
- Interactive review of key decision points
- Text explanations of AI reasoning
- Recommendations for improvement

### Transition
- Text summary of final score before detailed analysis
- Interactive elements to replay key decisions

## Error Handling

### Game Logic Errors
- Validate all moves against game rules
- Prevent illegal card selections or captures
- Clear error messaging for invalid actions

### Technical Errors
- Console logging for debugging
- Inform the user that something went wrong 

## Testing Plan

### Acceptance Tests
- Treat every user story as a transaction: user action/visible reactions
- Test at the level of the application services: given user action, what is the returned UI model

### Unit Tests
- Test all core game rules in isolation
- Verify scoring logic for different scenarios
- Validate AI decision-making process

### User Testing
- Verify UI usability for selecting cards and combinations
- Test post-game analysis clarity and helpfulness
- Validate AI strength and learning value

## Future Enhancements (Post-MVP)
- Statistics tracking across multiple games
- Multiple language support
- Accessibility features
- Additional difficulty levels for AI

## Development Priorities
1. Basic UI rendering and interaction
2. Core game engine and rules implementation
3. AI opponent implementation
4. Post-game analysis features