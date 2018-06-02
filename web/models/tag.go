package models

import (
	"net/url"
	"time"

	"github.com/luizbranco/srs/web"
)

type Tag struct {
	MetaID        web.ID    `db:"id"`
	MetaVersion   int       `db:"version"`
	MetaCreatedAt time.Time `db:"created_at"`
	MetaUpdatedAt time.Time `db:"updated_at"`

	DeckID web.ID `db:"deck_id"`
	Name   string `db:"name"`
}

func NewTagFromForm(deckID web.ID, form url.Values) (*Tag, error) {
	t := &Tag{
		DeckID: deckID,
		Name:   form.Get("name"),
	}
	return t, nil
}

func (t *Tag) ID() web.ID {
	return t.MetaID
}

func (t *Tag) SetID(id web.ID) {
	t.MetaID = id
}

func (t *Tag) Type() string {
	return "tag"
}

type CardTag struct {
	MetaID web.ID `db:"id"`
	CardID web.ID `db:"card_id"`
	TagID  web.ID `db:"tag_id"`
}

func (ct *CardTag) ID() web.ID {
	return ct.MetaID
}

func (ct *CardTag) SetID(id web.ID) {
	ct.MetaID = id
}

func (ct *CardTag) Type() string {
	return "card_tag"
}
