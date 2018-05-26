package web

import "github.com/gosimple/slug"

type Deck struct {
	ID          uint   `db:"id"`
	Slug        string `db:"slug"`
	Name        string `db:"name"`
	Description string `db:"description"`
}

func (d *Deck) Type() string {
	return "decks"
}

func (d *Deck) SetID(id uint) {
	d.ID = id
}

func (d *Deck) GenerateSlug() error {
	d.Slug = slug.Make(d.Name)
	return nil
}

type Decks []Deck

func (d *Decks) NewRecord() Record {
	return &Deck{}
}

func (d *Decks) Append(r Record) {
	deck := r.(*Deck)
	*d = append(*d, *deck)
}
