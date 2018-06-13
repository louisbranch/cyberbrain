package html

import (
	"github.com/pkg/errors"
	"gitlab.com/luizbranco/cyberbrain/primitives"
	"gitlab.com/luizbranco/cyberbrain/web"
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

	CreateCardPath     string
	CreateTagPath      string
	CreatePracticePath string
}

type Card struct {
	ID          string
	ImageURL    string
	SoundURL    string
	Definitions []string

	Path string

	Deck *Deck
	Tags []*Tag
}

type Tag struct {
	ID   string
	Name string

	Path string

	Deck  *Deck
	Cards []*Card
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

func RenderDeck(ub web.URLBuilder, d primitives.Deck, cards []primitives.Card,
	tags []primitives.Tag) (*Deck, error) {

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

	cp, err := ub.Path("NEW", &primitives.Card{}, d)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build deck new card path")
	}

	dr.NewCardPath = cp

	cp, err = ub.Path("CREATE", &primitives.Card{}, d)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build deck create card path")
	}

	dr.CreateCardPath = cp

	tp, err := ub.Path("NEW", &primitives.Tag{}, d)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build deck new tag path")
	}

	dr.NewTagPath = tp

	tp, err = ub.Path("CREATE", &primitives.Tag{}, d)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build deck create tag path")
	}

	dr.CreateTagPath = tp

	pp, err := ub.Path("NEW", &primitives.Practice{}, d)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build deck new practice path")
	}

	dr.NewPracticePath = pp

	pp, err = ub.Path("CREATE", &primitives.Practice{}, d)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build deck create practice path")
	}

	dr.CreatePracticePath = pp

	for _, c := range cards {
		cr, err := RenderCard(ub, d, nil, c, nil, false)
		if err != nil {
			return nil, errors.Wrap(err, "failed to render deck card")
		}

		dr.Cards = append(dr.Cards, cr)
	}

	for _, t := range tags {
		tr, err := RenderTag(ub, d, t, nil, false)
		if err != nil {
			return nil, errors.Wrap(err, "failed to render deck tag")
		}

		dr.Tags = append(dr.Tags, tr)
	}

	return dr, nil
}

func RenderCard(ub web.URLBuilder, d primitives.Deck, deckTags []primitives.Tag,
	c primitives.Card, cardTags []primitives.Tag, recursive bool) (*Card, error) {

	cr := &Card{
		ImageURL:    c.ImageURL,
		SoundURL:    c.SoundURL,
		Definitions: c.Definitions,
	}

	id, err := ub.EncodeID(c.ID())
	if err != nil {
		return nil, errors.Wrap(err, "failed to encode card id")
	}

	cr.ID = id

	p, err := ub.Path("SHOW", c, d)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build card path")
	}

	cr.Path = p

	for _, t := range cardTags {
		tr, err := RenderTag(ub, d, t, nil, false)
		if err != nil {
			return nil, errors.Wrap(err, "failed to render card tag")
		}

		cr.Tags = append(cr.Tags, tr)
	}

	if recursive {
		dr, err := RenderDeck(ub, d, nil, deckTags)
		if err != nil {
			return nil, errors.Wrap(err, "failed to render card deck")
		}
		cr.Deck = dr
	}

	return cr, nil
}

func RenderTag(ub web.URLBuilder, d primitives.Deck, t primitives.Tag,
	cards []primitives.Card, recursive bool) (*Tag, error) {

	tr := &Tag{
		Name: t.Name,
	}

	id, err := ub.EncodeID(t.ID())
	if err != nil {
		return nil, errors.Wrap(err, "failed to encode tag id")
	}

	tr.ID = id

	p, err := ub.Path("SHOW", t, d)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build tag path")
	}

	tr.Path = p

	for _, c := range cards {
		cr, err := RenderCard(ub, d, nil, c, nil, false)
		if err != nil {
			return nil, errors.Wrap(err, "failed to render tag card")
		}

		tr.Cards = append(tr.Cards, cr)
	}

	if recursive {
		dr, err := RenderDeck(ub, d, nil, nil)
		if err != nil {
			return nil, errors.Wrap(err, "failed to render tag deck")
		}
		tr.Deck = dr
	}

	return tr, nil
}

func RenderPractice(ub web.URLBuilder, d primitives.Deck, p primitives.Practice,
	recursive bool) (*Practice, error) {

	pr := &Practice{
		Done: p.Done,
	}

	if p.Done {
		pr.State = "Finished"
	} else {
		pr.State = "In Progress"
	}

	path, err := ub.Path("SHOW", p, d)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build practice path")
	}

	pr.Path = path

	rp, err := ub.Path("NEW", &primitives.Round{}, p, d)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build continue practice path")
	}

	pr.NewRoundPath = rp

	if recursive {
		dr, err := RenderDeck(ub, d, nil, nil)
		if err != nil {
			return nil, errors.Wrap(err, "failed to render practice deck")
		}
		pr.Deck = dr
	}

	return pr, nil
}

func RenderRound(ub web.URLBuilder, d primitives.Deck, r primitives.Round,
	p primitives.Practice) (*Round, error) {

	rr := &Round{
		Answer:      r.Answer,
		Guess:       r.Guess,
		Done:        r.Done,
		Correct:     r.Correct,
		PromptImage: r.Prompt,
	}

	pr, err := RenderPractice(ub, d, p, true)
	if err != nil {
		return nil, errors.Wrap(err, "failed to render round practice")
	}
	rr.Practice = pr

	path, err := ub.Path("SHOW", r, d)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build round path")
	}

	rr.Path = path

	return rr, nil
}
