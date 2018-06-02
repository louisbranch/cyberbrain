package models

import (
	"github.com/luizbranco/srs/web"
	"github.com/pkg/errors"
)

type DeckRendered struct {
	ID          string
	Name        string
	Description string
	ImageURL    string
	CardFields  []string
	Tags        []*TagRendered
	Cards       []*CardRendered

	Path            string
	NewCardPath     string
	NewTagPath      string
	NewPracticePath string
}

type CardRendered struct {
	ID          string
	ImageURLs   []string
	AudioURLs   []string
	Definitions []string
	Tags        []*TagRendered

	Path string
}

type TagRendered struct {
	ID   string
	Name string

	Path string
}

func (d *Deck) Render(ub web.URLBuilder) (*DeckRendered, error) {
	dr := &DeckRendered{
		Name:        d.Name,
		Description: d.Description,
		ImageURL:    d.ImageURL,
		CardFields:  d.CardFields,
	}

	id, err := ub.EncodeID(d.ID())
	if err != nil {
		return nil, errors.Wrap(err, "failed to encode deck id")
	}

	dr.ID = id

	p, err := ub.Path("SHOW", d)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build deck path")
	}

	dr.Path = p

	cp, err := ub.Path("NEW", &Card{}, d)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build deck new card path")
	}

	dr.NewCardPath = cp

	tp, err := ub.Path("NEW", &Tag{}, d)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build deck new tag path")
	}

	dr.NewTagPath = tp

	pp, err := ub.Path("NEW", &Practice{}, d)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build deck new practice path")
	}

	dr.NewPracticePath = pp

	for _, c := range d.Cards {
		cr, err := c.Render(ub)
		if err != nil {
			return nil, errors.Wrap(err, "failed to render deck card")
		}

		dr.Cards = append(dr.Cards, cr)
	}

	for _, t := range d.Tags {
		tr, err := t.Render(ub)
		if err != nil {
			return nil, errors.Wrap(err, "failed to render deck tag")
		}

		dr.Tags = append(dr.Tags, tr)
	}

	return dr, nil
}

func (c *Card) Render(ub web.URLBuilder) (*CardRendered, error) {
	cr := &CardRendered{
		ImageURLs:   c.ImageURLs,
		AudioURLs:   c.AudioURLs,
		Definitions: c.Definitions,
	}

	id, err := ub.EncodeID(c.ID())
	if err != nil {
		return nil, errors.Wrap(err, "failed to encode card id")
	}

	cr.ID = id

	p, err := ub.Path("SHOW", c)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build card path")
	}

	cr.Path = p

	for _, t := range c.Tags {
		tr, err := t.Render(ub)
		if err != nil {
			return nil, errors.Wrap(err, "failed to render card tag")
		}

		cr.Tags = append(cr.Tags, tr)
	}

	return cr, nil
}

func (t *Tag) Render(ub web.URLBuilder) (*TagRendered, error) {
	tr := &TagRendered{
		Name: t.Name,
	}

	id, err := ub.EncodeID(t.ID())
	if err != nil {
		return nil, errors.Wrap(err, "failed to encode tag id")
	}

	tr.ID = id

	p, err := ub.Path("SHOW", t)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build tag path")
	}

	tr.Path = p

	return tr, nil
}
