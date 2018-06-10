package finder

import (
	"fmt"
	"net/http"
	"sort"

	"gitlab.com/luizbranco/srs/db"
	"gitlab.com/luizbranco/srs/primitives"
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

func Deck(conn primitives.Database, ub web.URLBuilder, i identifier,
	opt option) (*primitives.Deck, []primitives.Card, []primitives.Tag, error) {

	id, err := parseID(ub, i)
	if err != nil {
		return nil, nil, nil, err
	}

	deck, err := db.FindDeck(conn, id)
	if err != nil {
		return nil, nil, nil, response.WrapError(err, http.StatusBadRequest, "wrong deck id")
	}

	var cards []primitives.Card
	var tags []primitives.Tag

	if opt&WithCards > 0 {
		cards, err = db.FindCardsByDeck(conn, id)
		if err != nil {
			return nil, nil, nil, response.WrapError(err, http.StatusInternalServerError, "failed to find deck cards")
		}
	}

	if opt&WithTags > 0 {
		tags, err = db.FindTags(conn, id)
		if err != nil {
			return nil, nil, nil, response.WrapError(err, http.StatusInternalServerError, "failed to find deck tags")
		}

		sort.Slice(tags, func(i, j int) bool {
			return tags[i].Name < tags[j].Name
		})
	}

	return deck, cards, tags, nil
}

func Practice(conn primitives.Database, ub web.URLBuilder, i identifier) (*primitives.Practice, error) {
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

func Card(conn primitives.Database, ub web.URLBuilder, i identifier) (*primitives.Card,
	[]primitives.Tag, error) {

	id, err := parseID(ub, i)
	if err != nil {
		return nil, nil, err
	}

	card, err := db.FindCard(conn, id)
	if err != nil {
		return nil, nil, response.WrapError(err, http.StatusBadRequest, "wrong card id")
	}

	tags, err := db.FindTagsByCard(conn, card.ID())
	if err != nil {
		return nil, nil, response.WrapError(err, http.StatusInternalServerError, "failed to find card tags")
	}

	sort.Slice(tags, func(i, j int) bool {
		return tags[i].Name < tags[j].Name
	})

	return card, tags, nil
}

func Tag(conn primitives.Database, ub web.URLBuilder, i identifier,
	opt option) (*primitives.Tag,

	[]primitives.Card, error) {

	id, err := parseID(ub, i)
	if err != nil {
		return nil, nil, err
	}

	tag, err := db.FindTag(conn, id)
	if err != nil {
		return nil, nil, response.WrapError(err, http.StatusBadRequest, "wrong tag id")
	}

	var cards []primitives.Card

	if opt&WithCards > 0 {
		cards, err = db.FindCardsByTag(conn, tag.ID())
		if err != nil {
			return nil, nil, response.WrapError(err, http.StatusInternalServerError,
				"failed to find tag cards")
		}
	}

	return tag, cards, nil
}

func parseID(ub web.URLBuilder, i identifier) (primitives.ID, error) {
	var blank primitives.ID

	switch i.(type) {
	case primitives.ID:
		return i.(primitives.ID), nil
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
