package models

import (
	"net/url"

	"github.com/luizbranco/srs/web"
)

type Deck struct {
	ID          web.ID `db:"id"`
	Slug        string `db:"slug"`
	Name        string `db:"name"`
	Description string `db:"description"`
	ImageURL    string `db:"image_url"`
	Cards       []Card
	Tags        []Tag
}

func NewDeck() *Deck {
	return &Deck{Slug: NewSlug()}
}

func NewDeckFromForm(form url.Values) (*Deck, error) {
	d := NewDeck()

	d.Name = form.Get("name")
	d.Description = form.Get("description")
	d.ImageURL = form.Get("image_url")

	return d, nil
}

func (d *Deck) SetID(id web.ID) {
	d.ID = id
}

func (d *Deck) Type() string {
	return "deck"
}
