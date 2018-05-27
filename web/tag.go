package web

import "github.com/gosimple/slug"

type Tag struct {
	ID     uint   `db:"id"`
	DeckID uint   `db:"deck_id"`
	Slug   string `db:"slug"`
	Name   string `db:"name"`
}

func (t *Tag) Type() string {
	return "tags"
}

func (t *Tag) SetID(id uint) {
	t.ID = id
}

func (t *Tag) GenerateSlug() error {
	t.Slug = slug.Make(t.Name)
	return nil
}

type Tags []Tag

func (t *Tags) NewRecord() Record {
	return &Tag{}
}

func (t *Tags) Append(r Record) {
	tag := r.(*Tag)
	*t = append(*t, *tag)
}
