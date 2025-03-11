2025-03-10 Work on the card capture story.

Start at 170 requests.
I had to remove the bit about scoring, it was premature.

This prompt worked well:

Master Yoda, we need to implement the following user story:
<Story>
Please inspect the codebase and formulate a clever plan.  Explain the plan step by step to me, then review your own plan for flaws.  Also ask me for any missing information you may require.  Do not implement yet!  Let's reason first.

Captain Matt, before I proceed with implementation, I need to clarify a few points:
Does a card capture only one other card with the same value, or can it capture multiple cards whose sum equals the card's value?
How is scoring calculated in Scopa? What points are awarded for different actions?
When does the game end? After a certain number of rounds or points?
Should the AI make its move automatically after the player captures, or should there be a button to trigger the AI's turn?

This is a thin slice of functionality, it's not meant to be complete! Let's stick to implementing <the story>.  We need not worry now about <other>

Let me clarify: this is a thin slice of functionality, it's not meant to be complete! Let's stick to implementing the capture.  The basic idea is: user clicks on a card in their hand, then clicks on a card on the table, and if the two cards have the same rank, you have a capture!  We need not worry now about the turn switching, or when the game ends.

Please tell me if now all is clear, Master Yoda.  If no further questions you have, please start with implementing updated tests for all the parts of the code that are involved, and then let's stop and run the tests and see them fail, TDD style.


At 178 I start the implementation of the Random AI Player.  We end at 184.

2025-03-11 Reveiw

TODO

- I find that one of the tests is not working
- No need to test for GET in handlers
- Name and Value fields of Card should be methods
- Handler tests are too long and involved
- handlers should not all share the same struct?
- AI turn should be invoked by a separate action

