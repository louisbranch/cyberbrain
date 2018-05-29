package web

import "fmt"

type Card struct {
	ID            uint   `db:"id"`
	DeckID        uint   `db:"deck_id"`
	ImageURL      string `db:"image_url"`
	AudioURL      string `db:"audio_url"`
	Definition    string `db:"definition"`
	AltDefinition string `db:"alt_definition"`
	Pronunciation string `db:"pronunciation"`
	Tags          []Tag
}

func (c *Card) Type() string {
	return "cards"
}

func (c *Card) SetID(id uint) {
	c.ID = id
}

func (c *Card) GenerateSlug() error {
	return nil
}

func (c *Card) Slug() string {
	return fmt.Sprintf("%d", c.ID)
}
