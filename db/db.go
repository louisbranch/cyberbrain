package db

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"gitlab.com/luizbranco/cyberbrain/primitives"
)

var ErrNotEnoughCards = errors.New("not enough cards")

type LessOrEqual struct {
	Time time.Time
}

type GreaterOrEqual struct {
	Time time.Time
}

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

func FindCardsByDeck(db primitives.Database, deckID primitives.ID, nsfw bool) ([]primitives.Card, error) {
	q := newCardQuery()
	q.where["deck_id"] = deckID
	if !nsfw {
		q.where["nsfw"] = false
	}

	q.sortBy["updated_at"] = "DESC"

	rs, err := db.Query(q)
	if err != nil {
		return nil, err
	}

	return castCards(rs)
}

func FindCardsByTag(db primitives.Database, tagID primitives.ID) ([]primitives.Card, error) {
	raw := `SELECT c.* FROM cards c
	LEFT JOIN card_tags ct ON c.id = ct.card_id
	WHERE ct.tag_id = ` + tagID.String() + " ORDER BY c.updated_at DESC;"

	q := newCardQuery()
	q.raw = raw

	rs, err := db.QueryRaw(q)
	if err != nil {
		return nil, err
	}

	return castCards(rs)
}

func FindTagsByCard(db primitives.Database, cardID primitives.ID) ([]primitives.Tag, error) {
	raw := `SELECT t.* FROM tags t
	LEFT JOIN card_tags ct ON t.id = ct.tag_id
	WHERE ct.card_id = ` + cardID.String() + ";"

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

func CountCardsScheduled(db primitives.Database, deckID primitives.ID) (int, error) {
	q := newCardScheduleQuery()
	q.where["deck_id"] = deckID
	q.where["next_date"] = LessOrEqual{time.Now().UTC()}

	n, err := db.Count(q)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to count cards scheduled for deck %d", deckID)
	}

	return n, nil
}

func FindNextCardScheduled(db primitives.Database, deckID primitives.ID,
	nsfw bool) (*primitives.Card, error) {

	now := time.Now().UTC().Format(time.RFC3339)

	nsfwWhere := ""
	if !nsfw {
		nsfwWhere = "AND c.nsfw = false"
	}

	raw := fmt.Sprintf(`SELECT c.* FROM cards c
	RIGHT JOIN card_schedules cd ON c.id = cd.card_id
	WHERE c.deck_id = %s
	AND cd.next_date <= '%s'::date
	%s
	ORDER BY random()
	LIMIT 1;
	`, deckID, now, nsfwWhere)

	q := newCardQuery()
	q.raw = raw

	rs, err := db.QueryRaw(q)
	if err != nil {
		return nil, err
	}

	cards, err := castCards(rs)
	if err != nil {
		return nil, err
	}

	if len(cards) == 0 {
		return nil, errors.New("no next card schedule found")
	}

	return &cards[0], nil
}

func FindCardSchedule(db primitives.Database, cardID primitives.ID) (*primitives.CardSchedule, error) {
	q := newCardScheduleQuery()
	q.where["card_id"] = cardID

	r, err := db.Get(q)
	if err != nil {
		return nil, err
	}

	schedule, ok := r.(*primitives.CardSchedule)
	if !ok {
		return nil, errors.Errorf("invalid record type %T", r)
	}

	return schedule, nil
}

func FindCardReview(db primitives.Database, id primitives.ID) (
	*primitives.CardReview, error) {

	q := newCardReviewQuery()
	q.where["id"] = id

	r, err := db.Get(q)
	if err != nil {
		return nil, err
	}

	review, ok := r.(*primitives.CardReview)
	if !ok {
		return nil, errors.Errorf("invalid record type %T", r)
	}

	return review, nil
}
