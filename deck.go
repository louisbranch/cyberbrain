package srs

import (
	"time"
)

type Deck struct {
	MetaID        ID        `db:"id"`
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

func (d Deck) ID() ID {
	return d.MetaID
}

func (d *Deck) SetID(id ID) {
	d.MetaID = id
}

func (d Deck) Type() string {
	return "deck"
}
