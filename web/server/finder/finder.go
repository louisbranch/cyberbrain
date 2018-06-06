package finder

import (
	"fmt"
	"net/http"
	"sort"

	"gitlab.com/luizbranco/srs"
	"gitlab.com/luizbranco/srs/db"
	"gitlab.com/luizbranco/srs/web"
	"gitlab.com/luizbranco/srs/web/server/response"
)

type identifier interface{}

type option int

const (
	NoOption option = 0
	WithTags        = 1 << iota
	WithCards
)

func Deck(conn srs.Database, ub web.URLBuilder, i identifier, opt option) (*srs.Deck, error) {
	id, err := parseID(ub, i)
	if err != nil {
		return nil, err
	}

	deck, err := db.FindDeck(conn, id)
	if err != nil {
		return nil, response.WrapError(err, http.StatusBadRequest, "wrong deck id")
	}

	if opt&WithCards > 0 {
		cards, err := db.FindCardsByDeck(conn, id)
		if err != nil {
			return nil, response.WrapError(err, http.StatusInternalServerError, "failed to find deck cards")
		}

		deck.Cards = cards
	}

	if opt&WithTags > 0 {
		tags, err := db.FindTags(conn, id)
		if err != nil {
			return nil, response.WrapError(err, http.StatusInternalServerError, "failed to find deck tags")
		}

		sort.Slice(tags, func(i, j int) bool {
			return tags[i].Name < tags[j].Name
		})

		deck.Tags = tags
	}

	return deck, nil
}

func Practice(conn srs.Database, ub web.URLBuilder, i identifier) (*srs.Practice, error) {
	id, err := parseID(ub, i)
	if err != nil {
		return nil, err
	}

	practice, err := db.FindPractice(conn, id)
	if err != nil {
		return nil, response.WrapError(err, http.StatusBadRequest, "wrong practice id")
	}

	return practice, nil
}

func Card(conn srs.Database, ub web.URLBuilder, i identifier) (*srs.Card, error) {
	id, err := parseID(ub, i)
	if err != nil {
		return nil, err
	}

	card, err := db.FindCard(conn, id)
	if err != nil {
		return nil, response.WrapError(err, http.StatusBadRequest, "wrong card id")
	}

	deck, err := db.FindDeck(conn, card.DeckID)
	if err != nil {
		return nil, response.WrapError(err, http.StatusInternalServerError, "failed to find card deck")
	}

	tags, err := db.FindTagsByCard(conn, card.ID())
	if err != nil {
		return nil, response.WrapError(err, http.StatusInternalServerError, "failed to find card tags")
	}

	sort.Slice(tags, func(i, j int) bool {
		return tags[i].Name < tags[j].Name
	})

	card.Deck = deck
	card.Tags = tags

	return card, nil
}

func Tag(conn srs.Database, ub web.URLBuilder, i identifier) (*srs.Tag, error) {
	id, err := parseID(ub, i)
	if err != nil {
		return nil, err
	}

	tag, err := db.FindTag(conn, id)
	if err != nil {
		return nil, response.WrapError(err, http.StatusBadRequest, "wrong tag id")
	}

	deck, err := db.FindDeck(conn, tag.DeckID)
	if err != nil {
		return nil, response.WrapError(err, http.StatusInternalServerError, "failed to find tag deck")
	}

	cards, err := db.FindCardsByTag(conn, tag.ID())
	if err != nil {
		return nil, response.WrapError(err, http.StatusInternalServerError, "failed to find tag cards")
	}

	tag.Deck = deck
	tag.Cards = cards

	return tag, nil
}

func parseID(ub web.URLBuilder, i identifier) (srs.ID, error) {
	var blank srs.ID

	switch i.(type) {
	case srs.ID:
		return i.(srs.ID), nil
	case string:
		hash := i.(string)
		id, err := ub.ParseID(hash)
		if err != nil {
			msg := fmt.Sprintf("invalid id %q", hash)
			return blank, response.WrapError(err, http.StatusBadRequest, msg)
		}
		return id, nil
	default:
		msg := fmt.Sprintf("invalid id format %v", i)
		return blank, response.NewError(http.StatusInternalServerError, msg)
	}
}
