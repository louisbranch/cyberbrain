package web

import "fmt"

type Card struct {
	ID       uint   `db:"id"`
	DeckID   uint   `db:"deck_id"`
	Slug     string `db:"slug"`
	ImageURL string `db:"image_url"`
	AudioURL string `db:"audio_url"`
	Field1   string `db:"field_1"`
	Field2   string `db:"field_2"`
	Field3   string `db:"field_3"`
}

func (c *Card) Type() string {
	return "cards"
}

func (c *Card) SetID(id uint) {
	c.ID = id
}

func (c *Card) GenerateSlug() error {
	c.Slug = fmt.Sprintf("%d", c.ID)
	return nil
}

type Cards []Card

func (c *Cards) NewRecord() Record {
	return &Card{}
}

func (c *Cards) Append(r Record) {
	card := r.(*Card)
	*c = append(*c, *card)
}
