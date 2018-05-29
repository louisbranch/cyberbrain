package models

import (
	"net/url"

	"github.com/luizbranco/srs/web"
)

type Card struct {
	ID            web.ID `db:"id"`
	DeckID        web.ID `db:"deck_id"`
	Slug          string `db:"slug"`
	ImageURL      string `db:"image_url"`
	AudioURL      string `db:"audio_url"`
	Definition    string `db:"definition"`
	AltDefinition string `db:"alt_definition"`
	Pronunciation string `db:"pronunciation"`
	Tags          []Tag
}

func NewCard() *Card {
	return &Card{Slug: NewSlug()}
}

func NewCardFromForm(deckID web.ID, form url.Values) (*Card, error) {
	c := NewCard()
	c.DeckID = deckID
	c.ImageURL = form.Get("image_url")
	c.AudioURL = form.Get("audio_url")
	c.Definition = form.Get("definition")
	c.AltDefinition = form.Get("alt_definition")
	c.Pronunciation = form.Get("pronunciation")
	return c, nil
}

func (c *Card) SetID(id web.ID) {
	c.ID = id
}

func (c *Card) Type() string {
	return "card"
}
