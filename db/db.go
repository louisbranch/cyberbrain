package db

import (
	"strconv"

	"github.com/pkg/errors"
	"gitlab.com/luizbranco/srs/primitives"
)

func FindUser(db primitives.Database, id primitives.ID) (*primitives.User, error) {
	q := newUserQuery()
	q.where["id"] = id

	r, err := db.Get(q)
	if err != nil {
		return nil, err
	}

	user, ok := r.(*primitives.User)
	if !ok {
		return nil, errors.Errorf("invalid record type %T", r)
	}

	return user, nil
}

func FindUserByEmail(db primitives.Database, email string) (*primitives.User, error) {
	q := newUserQuery()
	q.where["email"] = email

	r, err := db.Get(q)
	if err != nil {
		return nil, err
	}

	user, ok := r.(*primitives.User)
	if !ok {
		return nil, errors.Errorf("invalid record type %T", r)
	}

	return user, nil
}

func FindDecks(db primitives.Database, userID primitives.ID) ([]primitives.Deck, error) {
	q := newDeckQuery()
	q.where["user_id"] = userID

	rs, err := db.Query(q)
	if err != nil {
		return nil, err
	}

	var decks []primitives.Deck

	for _, r := range rs {
		deck, ok := r.(*primitives.Deck)
		if !ok {
			return nil, errors.Errorf("invalid record type %T", r)
		}

		decks = append(decks, *deck)
	}

	return decks, nil
}

func FindDeck(db primitives.Database, id primitives.ID) (*primitives.Deck, error) {
	q := newDeckQuery()
	q.where["id"] = id

	r, err := db.Get(q)
	if err != nil {
		return nil, err
	}

	deck, ok := r.(*primitives.Deck)
	if !ok {
		return nil, errors.Errorf("invalid record type %T", r)
	}

	return deck, nil
}

func FindCard(db primitives.Database, id primitives.ID) (*primitives.Card, error) {
	q := newCardQuery()
	q.where["id"] = id

	r, err := db.Get(q)
	if err != nil {
		return nil, err
	}

	card, ok := r.(*primitives.Card)
	if !ok {
		return nil, errors.Errorf("invalid record type %T", r)
	}

	return card, nil
}

func FindCardsByDeck(db primitives.Database, deckID primitives.ID) ([]primitives.Card, error) {
	q := newCardQuery()
	q.where["deck_id"] = deckID
	q.sortBy["updated_at"] = "DESC"

	rs, err := db.Query(q)
	if err != nil {
		return nil, err
	}

	return castCards(rs)
}

func FindCardsByTag(db primitives.Database, tagID primitives.ID) ([]primitives.Card, error) {
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

func FindTagsByCard(db primitives.Database, cardID primitives.ID) ([]primitives.Tag, error) {
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

func FindTag(db primitives.Database, id primitives.ID) (*primitives.Tag, error) {
	q := newTagQuery()
	q.where["id"] = id

	r, err := db.Get(q)
	if err != nil {
		return nil, err
	}

	tag, ok := r.(*primitives.Tag)
	if !ok {
		return nil, errors.Errorf("invalid record type %T", r)
	}

	return tag, nil
}

func FindTags(db primitives.Database, deckID primitives.ID) ([]primitives.Tag, error) {
	q := newTagQuery()
	q.where["deck_id"] = deckID

	rs, err := db.Query(q)
	if err != nil {
		return nil, err
	}

	return castTags(rs)
}

func castTags(rs []primitives.Record) ([]primitives.Tag, error) {
	var tags []primitives.Tag

	for _, r := range rs {
		tag, ok := r.(*primitives.Tag)

		if !ok {
			return nil, errors.Errorf("invalid record type %T", r)
		}

		tags = append(tags, *tag)
	}

	return tags, nil
}

func FindPractice(db primitives.Database, id primitives.ID) (*primitives.Practice, error) {
	q := newPracticeQuery()
	q.where["id"] = id

	r, err := db.Get(q)
	if err != nil {
		return nil, err
	}

	p, ok := r.(*primitives.Practice)

	if !ok {
		return nil, errors.Errorf("invalid record type %T", r)
	}

	return p, nil
}

func FindRound(db primitives.Database, id primitives.ID) (*primitives.Round, error) {
	q := newRoundQuery()
	q.where["id"] = id

	r, err := db.Get(q)
	if err != nil {
		return nil, err
	}

	p, ok := r.(*primitives.Round)
	if !ok {
		return nil, errors.Errorf("invalid record type %T", r)
	}

	return p, nil
}

func FindRounds(db primitives.Database, practiceID primitives.ID) ([]primitives.Round, error) {
	q := newRoundQuery()
	q.where["practice_id"] = practiceID

	rs, err := db.Query(q)
	if err != nil {
		return nil, err
	}

	return castRounds(rs)
}

func CountRounds(db primitives.Database, practiceID primitives.ID) (int, error) {
	q := newRoundQuery()
	q.where["practice_id"] = practiceID

	return db.Count(q)
}

func RandomCard(db primitives.Database, deckID primitives.ID) (*primitives.Card, error) {
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

func castCards(rs []primitives.Record) ([]primitives.Card, error) {
	var cards []primitives.Card

	for _, r := range rs {
		card, ok := r.(*primitives.Card)

		if !ok {
			return nil, errors.Errorf("invalid record type %T", r)
		}

		cards = append(cards, *card)
	}

	return cards, nil
}

func castRounds(rs []primitives.Record) ([]primitives.Round, error) {
	var rounds []primitives.Round

	for _, r := range rs {
		round, ok := r.(*primitives.Round)

		if !ok {
			return nil, errors.Errorf("invalid record type %T", r)
		}

		rounds = append(rounds, *round)
	}

	return rounds, nil
}
