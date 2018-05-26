package srs

// Card represents a single flash card.
type Card interface {
	ID() string
	Definitions() Definitions
	Image() []byte
	Audio() []byte
}

// Deck represents a group of flash cards.
type Deck interface {
	ID() string
	Name() string
	Description() string
	Cards() []Card
	AddCard(Card) error
}

// Definitions is map with the type of a definition as key and the definition
// itself as a value.
type Definitions map[string]string
