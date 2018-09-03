package html

import (
	"github.com/pkg/errors"
	"gitlab.com/luizbranco/cyberbrain/primitives"
	"gitlab.com/luizbranco/cyberbrain/web"
)

type Deck struct {
	ID             string
	Name           string
	Description    string
	ImageURL       string
	Fields         []string
	CardsScheduled int
	Tags           []*Tag
	Cards          []*Card

	Path              string
	EditPath          string
	NewCardPath       string
	NewTagPath        string
	NewCardReviewPath string

	CreateCardPath string
	CreateTagPath  string
}

type Card struct {
	ID          string
	ImageURL    string
	SoundURL    string
	Caption     string
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
	dr.EditPath = p + "/edit"

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

	rp, err := ub.Path("NEW", &primitives.CardReview{}, d)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build deck new review path")
	}

	dr.NewCardReviewPath = rp

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

func RenderCard(ub web.URLBuilder, deck primitives.Deck, deckTags []primitives.Tag,
	card primitives.Card, cardTags []primitives.Tag, recursive bool) (*Card, error) {

	cr := &Card{
		ImageURL:    card.ImageURL,
		SoundURL:    card.SoundURL,
		Caption:     card.Caption,
		Definitions: card.Definitions,
	}

	id, err := ub.EncodeID(card.ID())
	if err != nil {
		return nil, errors.Wrap(err, "failed to encode card id")
	}

	cr.ID = id

	p, err := ub.Path("SHOW", card, deck)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build card path")
	}

	cr.Path = p

	for _, t := range cardTags {
		tr, err := RenderTag(ub, deck, t, nil, false)
		if err != nil {
			return nil, errors.Wrap(err, "failed to render card tag")
		}

		cr.Tags = append(cr.Tags, tr)
	}

	if recursive {
		dr, err := RenderDeck(ub, deck, nil, deckTags)
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
