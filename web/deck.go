package web

import "github.com/gosimple/slug"

type Deck struct {
	ID          uint   `db:"id"`
	Slug        string `db:"slug"`
	Name        string `db:"name"`
	Description string `db:"description"`
	ImageURL    string `db:"image_url"`
	Cards       []Card
	Tags        []Tag
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
