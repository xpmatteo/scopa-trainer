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


At 178 I start the implementation of the Random AI Player.  We end at 184.  it took 6 requests.

2025-03-11 Review the code and refactor.

We start at 184 requests.

TODO

- ✅ I find that one of the tests is not working
- ✅ No need to test for GET in handlers
- ✅ Name and Value fields of Card should be methods
- ✅ Handler tests are too long and involved
- AI turn should be invoked by a separate action
- Last AI action should be displayed: cards captured, or new card on the table
- 

We end at 189 - it took 5 requests

## 2025-03-13

Refactored the tests to follow Esko Luontola's way explained in https://martinfowler.com/articles/tdd-html-templates.html

Asked Yoda to add missing tests, and an AI rule to apply the method.

The rule is good, the tests also but I think he added too many tests.  Also I'm not sure that all the new rules have the correct numeric prefix, but I'm too much in a hurry now to check.

Started at 189, now at 194 -- 5 requests.


## 2025-03-15

Started at 213
Improved the template tests to make them more readable and make failures easier to fix

At 214: implement the button to make the AI play its turn.  Had to insist on simplifying the template tests, we arrived at 221.

Now we implement a new story "Deal New Cards When Hand is Empty"

At 224 I ask it to refactor the two flags GameInProgress and IsPlayerTurn into a single GameStatus enumerative.  The refactor is partially successful, as it's still using multiple booleans in the UIModel, but let that be for now.

We are now at 229

Had an interesting conversation on why the go- rules are not applied!

I think there are many non-tabular tests

Next time: ask it to review the tests in the view of the testing rule


2025-03-17

Refactored the AI rule names

Implemented the real card images (2 requests)

At 5: implement simplified scoring

I'm testing the new scores.  The analysis workflow rule did not work, it did not wait for the answers to questions
