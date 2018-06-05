package finder

import (
	"fmt"
	"net/http"

	"gitlab.com/luizbranco/srs"
	"gitlab.com/luizbranco/srs/db"
	"gitlab.com/luizbranco/srs/web"
	"gitlab.com/luizbranco/srs/web/server/response"
)

type identifier interface{}

func DeckWithTags(conn srs.Database, ub web.URLBuilder, i identifier) (*srs.Deck, error) {
	id, err := parseID(ub, i)
	if err != nil {
		return nil, err
	}

	deck, err := db.FindDeck(conn, id)
	if err != nil {
		return nil, response.WrapError(err, http.StatusBadRequest, "wrong deck id")
	}

	tags, err := db.FindTags(conn, id)
	if err != nil {
		return nil, response.WrapError(err, http.StatusInternalServerError, "failed to find deck tags")
	}

	deck.Tags = tags

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

	card.Deck = deck
	card.Tags = tags

	return card, nil
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
