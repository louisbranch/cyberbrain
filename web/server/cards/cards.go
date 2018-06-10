package cards

import (
	"context"
	"fmt"
	"net/http"

	"gitlab.com/luizbranco/srs/db"
	"gitlab.com/luizbranco/srs/primitives"
	"gitlab.com/luizbranco/srs/web"
	"gitlab.com/luizbranco/srs/web/html"
	"gitlab.com/luizbranco/srs/web/server/finder"
	"gitlab.com/luizbranco/srs/web/server/response"
)

func Index() response.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) response.Responder {
		return response.Redirect{Path: "/decks/", Code: http.StatusFound}
	}
}

func New(conn primitives.Database, ub web.URLBuilder) response.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) response.Responder {

		query := r.URL.Query()
		hash := query.Get("deck")

		deck, err := finder.Deck(conn, ub, hash, finder.WithTags)
		if err != nil {
			return err.(response.Error)
		}

		content, err := html.RenderDeck(*deck, ub)
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

		hash := r.Form.Get("deck")

		deck, err := finder.Deck(conn, ub, hash, finder.NoOption)
		if err != nil {
			return err.(response.Error)
		}

		card, err := html.NewCardFromForm(*deck, r.Form)
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

		card, err := finder.Card(conn, ub, hash)
		if err != nil {
			return err.(response.Error)
		}

		content, err := html.RenderCard(*card, ub)
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

		card, err := finder.Card(conn, ub, hash)
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
