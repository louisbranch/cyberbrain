package primitives

import (
	"time"
)

type Tag struct {
	MetaID        ID        `db:"id"`
	MetaVersion   int       `db:"version"`
	MetaCreatedAt time.Time `db:"created_at"`
	MetaUpdatedAt time.Time `db:"updated_at"`

	DeckID ID     `db:"deck_id"`
	Name   string `db:"name"`
}

func (t Tag) ID() ID {
	return t.MetaID
}

func (t Tag) Type() string {
	return "tag"
}

func (t *Tag) SetID(id ID) {
	t.MetaID = id
}

func (t *Tag) SetVersion(v int) {
	t.MetaVersion = v
}

func (t *Tag) SetCreatedAt(ca time.Time) {
	t.MetaCreatedAt = ca
}

func (t *Tag) SetUpdatedAt(ua time.Time) {
	t.MetaUpdatedAt = ua
}

type CardTag struct {
	MetaID        ID        `db:"id"`
	MetaVersion   int       `db:"version"`
	MetaCreatedAt time.Time `db:"created_at"`
	MetaUpdatedAt time.Time `db:"updated_at"`

	CardID ID `db:"card_id"`
	TagID  ID `db:"tag_id"`
}

func (ct CardTag) ID() ID {
	return ct.MetaID
}

func (ct CardTag) Type() string {
	return "card_tag"
}

func (ct *CardTag) SetID(id ID) {
	ct.MetaID = id
}

func (ct *CardTag) SetVersion(v int) {
	ct.MetaVersion = v
}

func (ct *CardTag) SetCreatedAt(t time.Time) {
	ct.MetaCreatedAt = t
}

func (ct *CardTag) SetUpdatedAt(t time.Time) {
	ct.MetaUpdatedAt = t
}
