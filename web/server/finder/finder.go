package finder

import (
	"fmt"
	"net/http"
	"sort"

	"gitlab.com/luizbranco/cyberbrain/db"
	"gitlab.com/luizbranco/cyberbrain/primitives"
	"gitlab.com/luizbranco/cyberbrain/web"
	"gitlab.com/luizbranco/cyberbrain/web/server/response"
)

type identifier interface{}

type Option int

const (
	NoOption Option = 0
	WithTags Option = 1 << iota
	WithCards
	NSFW
)

func Deck(conn primitives.Database, ub web.URLBuilder, i identifier,
	opt Option) (*primitives.Deck, []primitives.Card, []primitives.Tag, error) {

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

	nsfw := false
	if opt&NSFW > 0 {
		nsfw = true
	}

	if opt&WithCards > 0 {
		cards, err = db.FindCardsByDeck(conn, id, nsfw)
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

func Card(conn primitives.Database, ub web.URLBuilder, i identifier, opt Option) (*primitives.Card,
	[]primitives.Tag, error) {

	id, err := parseID(ub, i)
	if err != nil {
		return nil, nil, err
	}

	card, err := db.FindCard(conn, id)
	if err != nil {
		return nil, nil, response.WrapError(err, http.StatusBadRequest, "wrong card id")
	}

	if opt&WithTags == 0 {
		return card, nil, nil
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
	opt Option) (*primitives.Tag,

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

func CardReview(conn primitives.Database, ub web.URLBuilder, i identifier) (
	*primitives.CardReview, error) {

	id, err := parseID(ub, i)
	if err != nil {
		return nil, err
	}

	review, err := db.FindCardReview(conn, id)
	if err != nil {
		return nil, response.WrapError(err, http.StatusBadRequest, "wrong review id")
	}

	return review, nil
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
