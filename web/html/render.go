package html

import (
	"github.com/luizbranco/srs"
	"github.com/luizbranco/srs/web"
	"github.com/pkg/errors"
)

type Deck struct {
	ID          string
	Name        string
	Description string
	ImageURL    string
	Fields      []string
	Tags        []*Tag
	Cards       []*Card

	Path            string
	NewCardPath     string
	NewTagPath      string
	NewPracticePath string
}

type Card struct {
	ID          string
	ImageURLs   []string
	AudioURLs   []string
	Definitions []string
	Tags        []*Tag

	Path string
}

type Tag struct {
	ID   string
	Name string

	Path string
}

type Practice struct {
	State string

	ContinuePath string

	// TODO progress

	Deck *Deck
}

func RenderDeck(d srs.Deck, ub web.URLBuilder) (*Deck, error) {
	dr := &Deck{
		Name:        d.Name,
		Description: d.Description,
		ImageURL:    d.ImageURL,
		Fields:      d.Fields,
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

	cp, err := ub.Path("NEW", &srs.Card{}, d)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build deck new card path")
	}

	dr.NewCardPath = cp

	tp, err := ub.Path("NEW", &srs.Tag{}, d)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build deck new tag path")
	}

	dr.NewTagPath = tp

	pp, err := ub.Path("NEW", &srs.Practice{}, d)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build deck new practice path")
	}

	dr.NewPracticePath = pp

	for _, c := range d.Cards {
		cr, err := RenderCard(c, ub)
		if err != nil {
			return nil, errors.Wrap(err, "failed to render deck card")
		}

		dr.Cards = append(dr.Cards, cr)
	}

	for _, t := range d.Tags {
		tr, err := RenderTag(t, ub)
		if err != nil {
			return nil, errors.Wrap(err, "failed to render deck tag")
		}

		dr.Tags = append(dr.Tags, tr)
	}

	return dr, nil
}

func RenderCard(c srs.Card, ub web.URLBuilder) (*Card, error) {
	cr := &Card{
		ImageURLs:   c.ImageURLs,
		AudioURLs:   c.SoundURLs,
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
		tr, err := RenderTag(t, ub)
		if err != nil {
			return nil, errors.Wrap(err, "failed to render card tag")
		}

		cr.Tags = append(cr.Tags, tr)
	}

	return cr, nil
}

func RenderTag(t srs.Tag, ub web.URLBuilder) (*Tag, error) {
	tr := &Tag{
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

func RenderPractice(p srs.Practice, ub web.URLBuilder) (*Practice, error) {
	pr := &Practice{}

	if p.Deck != nil {
		dr, err := RenderDeck(*p.Deck, ub)
		if err != nil {
			return nil, errors.Wrap(err, "failed to render deck practice")
		}
		pr.Deck = dr
	}

	if p.Done {
		pr.State = "Finished"
		return pr, nil
	}

	pr.State = "In Progress"

	path, err := ub.Path("NEW", &srs.PracticeRound{}, p)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build continue practice path")
	}

	pr.ContinuePath = path

	return pr, nil
}
