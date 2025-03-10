# Scopa Trainer

A browser-based implementation of the Scopa card game (two-player variant) that allows a human player to play against an AI opponent.

## Features

- Play the traditional Italian card game Scopa against an AI opponent
- Simple and intuitive user interface
- Automatic turn alternation between player and AI
- Card capturing mechanics based on matching ranks

## Random AI Player

The game includes a simple AI player that follows these rules:

1. Always plays the first card in its hand
2. If the card can capture a card of the same rank on the table, it captures it
3. If no capture is possible, the card is played to the table

The AI player automatically takes its turn after the human player completes their move.

## Game Rules (Simplified)

- Players take turns playing one card at a time
- A card can capture a card of the same rank from the table
- If no capture is possible, the card is played to the table
- The human player always plays first, followed by the AI

## Development

### Running the Tests

```
go test ./pkg/... -v
```

### Running the Application

```
go run cmd/server/main.go
```

Then open your browser to http://localhost:8080

## Future Enhancements

- Implement more complex capture mechanics (combinations of cards)
- Add scoring system
- Implement more sophisticated AI strategies
- Add post-game analysis 