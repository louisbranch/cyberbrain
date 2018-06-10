package primitives

import (
	"time"
)

type Deck struct {
	MetaID        ID        `db:"id"`
	MetaVersion   int       `db:"version"`
	MetaCreatedAt time.Time `db:"created_at"`
	MetaUpdatedAt time.Time `db:"updated_at"`

	UserID      ID       `db:"user_id"`
	Name        string   `db:"name"`
	Description string   `db:"description"`
	ImageURL    string   `db:"image_url"`
	Fields      []string `db:"fields"`
}

func (d Deck) ID() ID {
	return d.MetaID
}

func (d Deck) Type() string {
	return "deck"
}

func (d *Deck) SetID(id ID) {
	d.MetaID = id
}

func (d *Deck) SetVersion(v int) {
	d.MetaVersion = v
}

func (d *Deck) SetCreatedAt(t time.Time) {
	d.MetaCreatedAt = t
}

func (d *Deck) SetUpdatedAt(t time.Time) {
	d.MetaUpdatedAt = t
}
