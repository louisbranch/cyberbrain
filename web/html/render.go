package html

import (
	"github.com/pkg/errors"
	"gitlab.com/luizbranco/srs"
	"gitlab.com/luizbranco/srs/web"
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
	SoundURLs   []string
	Definitions []string
	Tags        []*Tag

	Path string

	Deck *Deck
}

type Tag struct {
	ID   string
	Name string

	Path string
}

type Practice struct {
	State string
	Done  bool

	Path         string
	NewRoundPath string

	// TODO progress

	Deck *Deck
}

type Round struct {
	PromptImage string
	Answer      string
	Guess       string
	Done        bool
	Correct     bool

	Practice *Practice

	Path string
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
		SoundURLs:   c.SoundURLs,
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

	if c.Deck != nil {
		dr, err := RenderDeck(*c.Deck, ub)
		if err != nil {
			return nil, errors.Wrap(err, "failed to render card deck")
		}
		cr.Deck = dr
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
	pr := &Practice{
		Done: p.Done,
	}

	if p.Deck != nil {
		dr, err := RenderDeck(*p.Deck, ub)
		if err != nil {
			return nil, errors.Wrap(err, "failed to render practice deck")
		}
		pr.Deck = dr
	}

	if p.Done {
		pr.State = "Finished"
	} else {
		pr.State = "In Progress"
	}

	path, err := ub.Path("SHOW", p)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build practice path")
	}

	pr.Path = path

	rp, err := ub.Path("NEW", &srs.Round{}, p)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build continue practice path")
	}

	pr.NewRoundPath = rp

	return pr, nil
}

func RenderRound(r srs.Round, ub web.URLBuilder) (*Round, error) {
	rr := &Round{
		Answer:      r.Answer,
		Guess:       r.Guess,
		Done:        r.Done,
		Correct:     r.Correct,
		PromptImage: r.Prompt,
	}

	if r.Practice != nil {
		p, err := RenderPractice(*r.Practice, ub)
		if err != nil {
			return nil, errors.Wrap(err, "failed to render round practice")
		}
		rr.Practice = p
	}

	p, err := ub.Path("SHOW", r)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build round path")
	}

	rr.Path = p

	return rr, nil
}
