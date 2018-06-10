package primitives

import (
	"time"
)

type Card struct {
	MetaID        ID        `db:"id"`
	MetaVersion   int       `db:"version"`
	MetaCreatedAt time.Time `db:"created_at"`
	MetaUpdatedAt time.Time `db:"updated_at"`

	DeckID      ID       `db:"deck_id"`
	Definitions []string `db:"definitions"`
	ImageURL    string   `db:"image_url"`
	SoundURL    string   `db:"sound_url"`
}

func (c Card) ID() ID {
	return c.MetaID
}

func (c Card) Type() string {
	return "card"
}

func (c *Card) SetID(id ID) {
	c.MetaID = id
}

func (d *Card) SetVersion(v int) {
	d.MetaVersion = v
}

func (d *Card) SetCreatedAt(t time.Time) {
	d.MetaCreatedAt = t
}

func (d *Card) SetUpdatedAt(t time.Time) {
	d.MetaUpdatedAt = t
}
