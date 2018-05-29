package server

import (
	"fmt"

	"github.com/luizbranco/srs/web"
	"github.com/pkg/errors"
)

func FindDecks(db web.Database) ([]web.Deck, error) {
	rec, err := db.Query(web.Condition{})
	if err != nil {
		return nil, err
	}

	var decks []web.Deck

	for _, r := range rec {
		deck, ok := r.(*web.Deck)

		if !ok {
			return nil, errors.Errorf("invalid record type %T", rec)
		}

		decks = append(decks, *deck)
	}

	return decks, nil
}

func FindDeckBySlug(db web.Database, slug string) (*web.Deck, error) {
	cond := web.Condition{
		Where: map[string]interface{}{
			"slug": slug,
		},
	}

	return findDeck(db, cond)
}

func FindDeckByID(db web.Database, id uint) (*web.Deck, error) {
	cond := web.Condition{
		Where: map[string]interface{}{
			"id": id,
		},
	}

	return findDeck(db, cond)
}

func findDeck(db web.Database, cond web.Condition) (*web.Deck, error) {
	rec, err := db.Get(cond)
	if err != nil {
		return nil, err
	}

	deck, ok := rec.(*web.Deck)

	if !ok {
		return nil, errors.Errorf("invalid record type %T", rec)
	}

	return deck, nil
}

func FindCardByID(db web.Database, id uint64) (*web.Card, error) {
	where := web.Condition{
		Where: map[string]interface{}{
			"id": id,
		},
	}

	rec, err := db.Get(where)
	if err != nil {
		return nil, err
	}

	card, ok := rec.(*web.Card)
	if !ok {
		return nil, errors.Errorf("invalid record type %T", rec)
	}

	return card, nil
}

func FindCardsByDeckID(db web.Database, id uint) ([]web.Card, error) {
	cond := web.Condition{
		Where: map[string]interface{}{
			"id": id,
		},
	}

	rec, err := db.Query(cond)
	if err != nil {
		return nil, err
	}

	return castCards(rec)
}

func FindTagsByCardID(db web.Database, id uint64) ([]web.Tag, error) {
	i := fmt.Sprintf("%d", id)

	raw := `SELECT t.id, t.deck_id, t.slug, name FROM tags t
	LEFT JOIN card_tags ct ON t.id = ct.tag_id
	WHERE ct.card_id = ` + i + ";"

	cond := web.Condition{
		Raw: raw,
	}

	rec, err := db.QueryRaw(cond)
	if err != nil {
		return nil, err
	}

	return castTags(rec)
}

func FindTagsByDeckID(db web.Database, id uint) ([]web.Tag, error) {
	cond := web.Condition{
		Where: map[string]interface{}{
			"id": id,
		},
	}

	rec, err := db.Query(cond)
	if err != nil {
		return nil, err
	}

	return castTags(rec)
}

func castTags(rec []web.Record) ([]web.Tag, error) {
	var tags []web.Tag

	for _, r := range rec {
		tag, ok := r.(*web.Tag)

		if !ok {
			return nil, errors.Errorf("invalid record type %T", rec)
		}

		tags = append(tags, *tag)
	}

	return tags, nil
}

func FindPracticeByID(db web.Database, id uint64) (*web.Practice, error) {
	cond := web.Condition{
		Where: map[string]interface{}{
			"id": id,
		},
	}

	rec, err := db.Get(cond)
	if err != nil {
		return nil, err
	}

	p, ok := rec.(*web.Practice)

	if !ok {
		return nil, errors.Errorf("invalid record type %T", rec)
	}

	return p, nil
}

func FindPracticeRound(db web.Database, pid uint, round int) (*web.PracticeRound, error) {
	cond := web.Condition{
		Where: map[string]interface{}{
			"round":       round,
			"practice_id": pid,
		},
	}

	rec, err := db.Get(cond)
	if err != nil {
		return nil, err
	}

	p, ok := rec.(*web.PracticeRound)

	if !ok {
		return nil, errors.Errorf("invalid record type %T", rec)
	}

	return p, nil
}

func CountPracticeRounds(db web.Database, id uint) (int, error) {
	r := web.PracticeRound{}

	cond := web.Condition{
		Record: &r,
		Where: map[string]interface{}{
			"id": id,
		},
	}

	return db.Count(cond)
}

func FindRandomCard(db web.Database, deckID uint) (*web.Card, error) {
	r := web.Card{}

	cond := web.Condition{
		Record: &r,
		Where: map[string]interface{}{
			"deck_id": deckID,
		},
	}

	rec, err := db.Random(cond, 1)
	if err != nil {
		return nil, err
	}

	cards, err := castCards(rec)
	if err != nil {
		return nil, err
	}

	if len(cards) == 0 {
		return nil, errors.New("not enough cards")
	}

	return &cards[0], nil
}

func castCards(rec []web.Record) ([]web.Card, error) {
	var cards []web.Card

	for _, r := range rec {
		card, ok := r.(*web.Card)

		if !ok {
			return nil, errors.Errorf("invalid record type %T", rec)
		}

		cards = append(cards, *card)
	}

	return cards, nil
}
