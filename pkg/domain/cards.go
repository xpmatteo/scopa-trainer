package domain

// Suit represents a card suit
type Suit string

// Rank represents a card rank
type Rank int

// Suit constants
const (
	Coppe   Suit = "Coppe"
	Denari  Suit = "Denari"
	Bastoni Suit = "Bastoni"
	Spade   Suit = "Spade"
)

// Rank constants and their names
const (
	Asso    Rank = 1
	Due     Rank = 2
	Tre     Rank = 3
	Quattro Rank = 4
	Cinque  Rank = 5
	Sei     Rank = 6
	Sette   Rank = 7
	Fante   Rank = 8
	Cavallo Rank = 9
	Re      Rank = 10
)

// rankNames maps ranks to their string representation
var rankNames = map[Rank]string{
	Asso:    "Asso",
	Due:     "Due",
	Tre:     "Tre",
	Quattro: "Quattro",
	Cinque:  "Cinque",
	Sei:     "Sei",
	Sette:   "Sette",
	Fante:   "Fante",
	Cavallo: "Cavallo",
	Re:      "Re",
}

// Card represents a playing card in the Italian deck
type Card struct {
	Suit Suit
	Rank Rank
}

// String returns the string representation of a Rank
func (r Rank) String() string {
	return rankNames[r]
}

// Value returns the numeric value of a Rank in the game
func (r Rank) Value() int {
	return int(r)
}

// Name returns the string representation of a Card's rank
func (c Card) Name() string {
	return c.Rank.String()
}

// Value returns the numeric value of a Card in the game
func (c Card) Value() int {
	return c.Rank.Value()
}

// AllSuits returns all available suits
func AllSuits() []Suit {
	return []Suit{Coppe, Denari, Bastoni, Spade}
}

// AllRanks returns all available ranks
func AllRanks() []Rank {
	return []Rank{Asso, Due, Tre, Quattro, Cinque, Sei, Sette, Fante, Cavallo, Re}
}
