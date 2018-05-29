package models

import (
	"net/url"

	"github.com/luizbranco/srs/web"
)

type Tag struct {
	ID     web.ID `db:"id"`
	DeckID web.ID `db:"deck_id"`
	Slug   string `db:"slug"`
	Name   string `db:"name"`
}

func NewTag() *Tag {
	return &Tag{Slug: NewSlug()}
}

func NewTagFromForm(deckID web.ID, form url.Values) (*Tag, error) {
	t := NewTag()
	t.DeckID = deckID
	t.Name = form.Get("name")
	return t, nil
}

func (t *Tag) SetID(id web.ID) {
	t.ID = id
}

func (t *Tag) Type() string {
	return "tag"
}

type CardTag struct {
	ID     web.ID `db:"id"`
	CardID web.ID `db:"card_id"`
	TagID  web.ID `db:"tag_id"`
}

func (ct *CardTag) SetID(id web.ID) {
	ct.ID = id
}

func (ct *CardTag) Type() string {
	return "card_tag"
}
