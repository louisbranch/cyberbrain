package cards

import (
	"context"
	"fmt"
	"net/http"
	"sort"

	"gitlab.com/luizbranco/srs/db"
	"gitlab.com/luizbranco/srs/primitives"
	"gitlab.com/luizbranco/srs/web"
	"gitlab.com/luizbranco/srs/web/html"
	"gitlab.com/luizbranco/srs/web/server/finder"
	"gitlab.com/luizbranco/srs/web/server/middlewares"
	"gitlab.com/luizbranco/srs/web/server/response"
)

func Index() response.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) response.Responder {
		return response.Redirect{Path: "/decks/", Code: http.StatusFound}
	}
}

func New(conn primitives.Database, ub web.URLBuilder) response.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) response.Responder {

		deck, _ := middlewares.CurrentDeck(ctx)

		tags, err := db.FindTags(conn, deck.ID())
		if err != nil {
			return response.WrapError(err, http.StatusInternalServerError, "failed to find deck tags")
		}

		sort.Slice(tags, func(i, j int) bool {
			return tags[i].Name < tags[j].Name
		})

		content, err := html.RenderDeck(ub, deck, nil, tags)
		if err != nil {
			return response.WrapError(err, http.StatusInternalServerError, "failed to render deck")
		}

		page := web.Page{
			Title:    "New Card",
			Partials: []string{"new_card"},
			Content:  content,
		}

		return response.NewContent(page)
	}
}

func Create(conn primitives.Database, ub web.URLBuilder) response.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) response.Responder {

		if err := r.ParseForm(); err != nil {
			return response.WrapError(err, http.StatusBadRequest, "invalid form")
		}

		deck, _ := middlewares.CurrentDeck(ctx)

		card, err := html.NewCardFromForm(deck, r.Form)
		if err != nil {
			return response.WrapError(err, http.StatusBadRequest, "invalid card form")
		}

		err = conn.Create(card)
		if err != nil {
			return response.WrapError(err, http.StatusInternalServerError, "failed to create card")
		}

		tags := r.Form["tags"]

		for _, tag := range tags {
			id, err := ub.ParseID(tag)
			if err != nil {
				msg := fmt.Sprintf("invalid tag id %s", tag)
				return response.WrapError(err, http.StatusBadRequest, msg)
			}

			ct := primitives.CardTag{
				CardID: card.MetaID,
				TagID:  id,
			}

			err = conn.Create(&ct)
			if err != nil {
				return response.WrapError(err, http.StatusInternalServerError, "failed to create card tag")
			}
		}

		path, err := ub.Path("SHOW", deck)
		if err != nil {
			return response.WrapError(err, http.StatusInternalServerError, "failed to generate deck path")
		}

		return response.Redirect{Path: path, Code: http.StatusFound}
	}
}

func Show(conn primitives.Database, ub web.URLBuilder, hash string) response.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) response.Responder {

		card, tags, err := finder.Card(conn, ub, hash)
		if err != nil {
			return err.(response.Error)
		}

		deck, _ := middlewares.CurrentDeck(ctx)

		content, err := html.RenderCard(ub, deck, *card, tags, true)
		if err != nil {
			return response.WrapError(err, http.StatusInternalServerError, "failed to render card")
		}

		page := web.Page{
			Title:    "Card",
			Partials: []string{"card"},
			Content:  content,
		}

		return response.NewContent(page)
	}
}

func Update(conn primitives.Database, ub web.URLBuilder, hash string) response.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) response.Responder {

		if err := r.ParseForm(); err != nil {
			return response.WrapError(err, http.StatusBadRequest, "invalid form")
		}

		card, _, err := finder.Card(conn, ub, hash)
		if err != nil {
			return err.(response.Error)
		}

		deck, err := db.FindDeck(conn, card.DeckID)
		if err != nil {
			return response.WrapError(err, http.StatusInternalServerError, "invalid card deck")
		}

		newCard, err := html.NewCardFromForm(*deck, r.Form)
		if err != nil {
			return response.WrapError(err, http.StatusBadRequest, "invalid card form")
		}

		card.ImageURLs = newCard.ImageURLs
		card.SoundURLs = newCard.SoundURLs
		card.Definitions = newCard.Definitions

		// TODO reassign card tags

		err = conn.Update(card)
		if err != nil {
			return response.WrapError(err, http.StatusBadRequest, "failed to update card")
		}

		path, err := ub.Path("SHOW", deck)
		if err != nil {
			return response.WrapError(err, http.StatusInternalServerError, "failed to generate deck path")
		}

		return response.Redirect{Path: path, Code: http.StatusFound}
	}
}
