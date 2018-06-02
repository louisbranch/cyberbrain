package models

import (
	"github.com/dchest/uniuri"
	"github.com/luizbranco/srs/web"
	"github.com/pkg/errors"
)

func FindDecks(db web.Database) ([]Deck, error) {
	q := newDeckQuery()

	rs, err := db.Query(q)
	if err != nil {
		return nil, err
	}

	var decks []Deck

	for _, r := range rs {
		deck, ok := r.(*Deck)
		if !ok {
			return nil, errors.Errorf("invalid record type %T", r)
		}

		decks = append(decks, *deck)
	}

	return decks, nil
}

func FindDeckBySlug(db web.Database, slug string) (*Deck, error) {
	q := newDeckQuery()
	q.where["slug"] = slug
	return findDeck(db, q)
}

func FindDeckByID(db web.Database, id web.ID) (*Deck, error) {
	q := newDeckQuery()
	q.where["id"] = id
	return findDeck(db, q)
}

func findDeck(db web.Database, q *query) (*Deck, error) {
	r, err := db.Get(q)
	if err != nil {
		return nil, err
	}

	deck, ok := r.(*Deck)
	if !ok {
		return nil, errors.Errorf("invalid record type %T", r)
	}

	return deck, nil
}

func FindCardByID(db web.Database, id web.ID) (*Card, error) {
	q := newCardQuery()
	q.where["id"] = id

	r, err := db.Get(q)
	if err != nil {
		return nil, err
	}

	card, ok := r.(*Card)
	if !ok {
		return nil, errors.Errorf("invalid record type %T", r)
	}

	return card, nil
}

func FindCardBySlug(db web.Database, slug string) (*Card, error) {
	q := newCardQuery()
	q.where["slug"] = slug

	r, err := db.Get(q)
	if err != nil {
		return nil, err
	}

	card, ok := r.(*Card)
	if !ok {
		return nil, errors.Errorf("invalid record type %T", r)
	}

	return card, nil
}

func FindCardsByDeckID(db web.Database, deckID web.ID) ([]Card, error) {
	q := newCardQuery()
	q.where["deck_id"] = deckID

	rs, err := db.Query(q)
	if err != nil {
		return nil, err
	}

	return castCards(rs)
}

func FindTagsByCardID(db web.Database, cardID web.ID) ([]Tag, error) {
	raw := `SELECT t.id, t.deck_id, t.slug, name FROM tags t
	LEFT JOIN card_tags ct ON t.id = ct.tag_id
	WHERE ct.card_id = ` + string(cardID) + ";"

	q := newTagQuery()
	q.raw = raw

	rs, err := db.QueryRaw(q)
	if err != nil {
		return nil, err
	}

	return castTags(rs)
}

func FindTagsByDeckID(db web.Database, id web.ID) ([]Tag, error) {
	q := newTagQuery()
	q.where["deck_id"] = id

	rs, err := db.Query(q)
	if err != nil {
		return nil, err
	}

	return castTags(rs)
}

func castTags(rs []web.Record) ([]Tag, error) {
	var tags []Tag

	for _, r := range rs {
		tag, ok := r.(*Tag)

		if !ok {
			return nil, errors.Errorf("invalid record type %T", r)
		}

		tags = append(tags, *tag)
	}

	return tags, nil
}

func FindPracticeBySlug(db web.Database, slug string) (*Practice, error) {
	q := newPracticeQuery()
	q.where["slug"] = slug

	r, err := db.Get(q)
	if err != nil {
		return nil, err
	}

	p, ok := r.(*Practice)

	if !ok {
		return nil, errors.Errorf("invalid record type %T", r)
	}

	return p, nil
}

func FindPracticeByID(db web.Database, id web.ID) (*Practice, error) {
	q := newPracticeQuery()
	q.where["id"] = id

	r, err := db.Get(q)
	if err != nil {
		return nil, err
	}

	p, ok := r.(*Practice)

	if !ok {
		return nil, errors.Errorf("invalid record type %T", r)
	}

	return p, nil
}

func FindPracticeRound(db web.Database, pid web.ID, round int) (*PracticeRound, error) {
	q := newPracticeRoundQuery()
	q.where["practice_id"] = pid
	q.where["round"] = round

	r, err := db.Get(q)
	if err != nil {
		return nil, err
	}

	p, ok := r.(*PracticeRound)
	if !ok {
		return nil, errors.Errorf("invalid record type %T", r)
	}

	return p, nil
}

func CountPracticeRounds(db web.Database, id web.ID) (int, error) {
	q := newPracticeRoundQuery()
	q.where["id"] = id

	return db.Count(q)
}

func FindRandomCard(db web.Database, deckID web.ID) (*Card, error) {
	q := newCardQuery()
	q.where["deck_id"] = deckID

	rs, err := db.Random(q, 1)
	if err != nil {
		return nil, err
	}

	cards, err := castCards(rs)
	if err != nil {
		return nil, err
	}

	if len(cards) == 0 {
		return nil, errors.New("not enough cards")
	}

	return &cards[0], nil
}

func castCards(rs []web.Record) ([]Card, error) {
	var cards []Card

	for _, r := range rs {
		card, ok := r.(*Card)

		if !ok {
			return nil, errors.Errorf("invalid record type %T", r)
		}

		cards = append(cards, *card)
	}

	return cards, nil
}

func NewSlug() string {
	return uniuri.New()
}
