package db

import (
	"strconv"

	"github.com/pkg/errors"
	"gitlab.com/luizbranco/srs"
)

func FindDecks(db srs.Database) ([]srs.Deck, error) {
	q := newDeckQuery()

	rs, err := db.Query(q)
	if err != nil {
		return nil, err
	}

	var decks []srs.Deck

	for _, r := range rs {
		deck, ok := r.(*srs.Deck)
		if !ok {
			return nil, errors.Errorf("invalid record type %T", r)
		}

		decks = append(decks, *deck)
	}

	return decks, nil
}

func FindDeck(db srs.Database, id srs.ID) (*srs.Deck, error) {
	q := newDeckQuery()
	q.where["id"] = id

	r, err := db.Get(q)
	if err != nil {
		return nil, err
	}

	deck, ok := r.(*srs.Deck)
	if !ok {
		return nil, errors.Errorf("invalid record type %T", r)
	}

	return deck, nil
}

func FindCard(db srs.Database, id srs.ID) (*srs.Card, error) {
	q := newCardQuery()
	q.where["id"] = id

	r, err := db.Get(q)
	if err != nil {
		return nil, err
	}

	card, ok := r.(*srs.Card)
	if !ok {
		return nil, errors.Errorf("invalid record type %T", r)
	}

	return card, nil
}

func FindCardsByDeck(db srs.Database, deckID srs.ID) ([]srs.Card, error) {
	q := newCardQuery()
	q.where["deck_id"] = deckID
	q.sortBy["updated_at"] = "DESC"

	rs, err := db.Query(q)
	if err != nil {
		return nil, err
	}

	return castCards(rs)
}

func FindCardsByTag(db srs.Database, tagID srs.ID) ([]srs.Card, error) {
	id := strconv.Itoa(int(tagID))

	raw := `SELECT c.* FROM cards c
	LEFT JOIN card_tags ct ON c.id = ct.card_id
	WHERE ct.tag_id = ` + id + " ORDER BY c.updated_at DESC;"

	q := newCardQuery()
	q.raw = raw

	rs, err := db.QueryRaw(q)
	if err != nil {
		return nil, err
	}

	return castCards(rs)
}

func FindTagsByCard(db srs.Database, cardID srs.ID) ([]srs.Tag, error) {
	id := strconv.Itoa(int(cardID))

	raw := `SELECT t.* FROM tags t
	LEFT JOIN card_tags ct ON t.id = ct.tag_id
	WHERE ct.card_id = ` + id + ";"

	q := newTagQuery()
	q.raw = raw

	rs, err := db.QueryRaw(q)
	if err != nil {
		return nil, err
	}

	return castTags(rs)
}

func FindTag(db srs.Database, id srs.ID) (*srs.Tag, error) {
	q := newTagQuery()
	q.where["id"] = id

	r, err := db.Get(q)
	if err != nil {
		return nil, err
	}

	tag, ok := r.(*srs.Tag)
	if !ok {
		return nil, errors.Errorf("invalid record type %T", r)
	}

	return tag, nil
}

func FindTags(db srs.Database, deckID srs.ID) ([]srs.Tag, error) {
	q := newTagQuery()
	q.where["deck_id"] = deckID

	rs, err := db.Query(q)
	if err != nil {
		return nil, err
	}

	return castTags(rs)
}

func castTags(rs []srs.Record) ([]srs.Tag, error) {
	var tags []srs.Tag

	for _, r := range rs {
		tag, ok := r.(*srs.Tag)

		if !ok {
			return nil, errors.Errorf("invalid record type %T", r)
		}

		tags = append(tags, *tag)
	}

	return tags, nil
}

func FindPractice(db srs.Database, id srs.ID) (*srs.Practice, error) {
	q := newPracticeQuery()
	q.where["id"] = id

	r, err := db.Get(q)
	if err != nil {
		return nil, err
	}

	p, ok := r.(*srs.Practice)

	if !ok {
		return nil, errors.Errorf("invalid record type %T", r)
	}

	return p, nil
}

func FindRound(db srs.Database, id srs.ID) (*srs.Round, error) {
	q := newRoundQuery()
	q.where["id"] = id

	r, err := db.Get(q)
	if err != nil {
		return nil, err
	}

	p, ok := r.(*srs.Round)
	if !ok {
		return nil, errors.Errorf("invalid record type %T", r)
	}

	return p, nil
}

func FindRounds(db srs.Database, practiceID srs.ID) ([]srs.Round, error) {
	q := newRoundQuery()
	q.where["practice_id"] = practiceID

	rs, err := db.Query(q)
	if err != nil {
		return nil, err
	}

	return castRounds(rs)
}

func CountRounds(db srs.Database, practiceID srs.ID) (int, error) {
	q := newRoundQuery()
	q.where["practice_id"] = practiceID

	return db.Count(q)
}

func RandomCard(db srs.Database, deckID srs.ID) (*srs.Card, error) {
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

func castCards(rs []srs.Record) ([]srs.Card, error) {
	var cards []srs.Card

	for _, r := range rs {
		card, ok := r.(*srs.Card)

		if !ok {
			return nil, errors.Errorf("invalid record type %T", r)
		}

		cards = append(cards, *card)
	}

	return cards, nil
}

func castRounds(rs []srs.Record) ([]srs.Round, error) {
	var rounds []srs.Round

	for _, r := range rs {
		round, ok := r.(*srs.Round)

		if !ok {
			return nil, errors.Errorf("invalid record type %T", r)
		}

		rounds = append(rounds, *round)
	}

	return rounds, nil
}
