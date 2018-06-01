package models

import (
	"net/url"
	"time"

	"github.com/luizbranco/srs/web"
)

type Deck struct {
	ID          web.ID    `db:"id"`
	Version     int       `db:"version"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	ImageURL    string    `db:"image_url"`
	CardFields  []string  `db:"card_fields"`
	Cards       []Card
	Tags        []Tag
}

func NewDeck() *Deck {
	return &Deck{}
}

func NewDeckFromForm(form url.Values) (*Deck, error) {
	d := NewDeck()

	d.Name = form.Get("name")
	d.Description = form.Get("description")
	d.ImageURL = form.Get("image_url")

	for _, cf := range form["card_fields"] {
		d.CardFields = append(d.CardFields, cf)
	}

	return d, nil
}

func (d *Deck) SetID(id web.ID) {
	d.ID = id
}

func (d *Deck) Type() string {
	return "deck"
}
