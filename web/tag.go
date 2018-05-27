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

func (t *Tags) Len() int {
	return len(*t)
}

func (t *Tags) Less(i, j int) bool {
	return (*t)[i].Slug < (*t)[j].Slug
}

func (t *Tags) Swap(i, j int) {
	(*t)[i], (*t)[j] = (*t)[j], (*t)[i]
}

type CardTag struct {
	ID     uint `db:"id"`
	CardID uint `db:"card_id"`
	TagID  uint `db:"tag_id"`
}

func (ct *CardTag) Type() string {
	return "card_tags"
}

func (ct *CardTag) SetID(id uint) {
	ct.ID = id
}

func (ct *CardTag) GenerateSlug() error {
	return nil
}
