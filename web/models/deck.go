package models

import (
	"net/url"
	"time"

	"github.com/luizbranco/srs/web"
)

type Deck struct {
	MetaID        web.ID    `db:"id"`
	MetaVersion   int       `db:"version"`
	MetaCreatedAt time.Time `db:"created_at"`
	MetaUpdatedAt time.Time `db:"updated_at"`

	Name        string   `db:"name"`
	Description string   `db:"description"`
	ImageURL    string   `db:"image_url"`
	CardFields  []string `db:"card_fields"`
	Cards       []Card
	Tags        []Tag
}

func NewDeckFromForm(form url.Values) (*Deck, error) {
	d := &Deck{}

	d.Name = form.Get("name")
	d.Description = form.Get("description")
	d.ImageURL = form.Get("image_url")

	for _, cf := range form["card_fields"] {
		d.CardFields = append(d.CardFields, cf)
	}

	return d, nil
}

func (d *Deck) ID() web.ID {
	return d.MetaID
}

func (d *Deck) SetID(id web.ID) {
	d.MetaID = id
}

func (d *Deck) Type() string {
	return "deck"
}
