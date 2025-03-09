
Let's implement this user story:

US2.2: Play Card from Hand
When it's the player's turn, they can select a card from their hand to play.
Acceptance Criteria:

1. Player can click on any card in their hand to select it
2. Selected card is visually highlighted
3. Playing on a different card selects it, and deselects previous one

Implementation notes:

To make this testable, we keep the logic in the GameService.  Selection highlight is a UI concern, so it's stored in the GameService state, not in the GameState.

  Avoid putting logic in JavaScript.  The click action is handled in the GameService.

