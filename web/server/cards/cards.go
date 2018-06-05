package cards

import (
	"fmt"
	"net/http"

	"gitlab.com/luizbranco/srs"
	"gitlab.com/luizbranco/srs/web"
	"gitlab.com/luizbranco/srs/web/html"
	"gitlab.com/luizbranco/srs/web/server/finder"
	"gitlab.com/luizbranco/srs/web/server/response"
)

func Index() response.Handler {
	return func(w http.ResponseWriter, r *http.Request) response.Responder {
		return response.Redirect{Path: "/decks/", Code: http.StatusFound}
	}
}

func New(conn srs.Database, ub web.URLBuilder) response.Handler {
	return func(w http.ResponseWriter, r *http.Request) response.Responder {

		query := r.URL.Query()
		hash := query.Get("deck")

		deck, err := finder.DeckWithTags(conn, ub, hash)
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

func Create(conn srs.Database, ub web.URLBuilder) response.Handler {
	return func(w http.ResponseWriter, r *http.Request) response.Responder {

		if err := r.ParseForm(); err != nil {
			return response.WrapError(err, http.StatusBadRequest, "invalid form")
		}

		hash := r.Form.Get("deck")

		deck, err := finder.DeckWithTags(conn, ub, hash)
		if err != nil {
			return err.(response.Error)
		}

		card, err := html.NewCardFromForm(*deck, r.Form)
		if err != nil {
			return response.WrapError(err, http.StatusBadRequest, "invalid card form")
		}

		err = conn.Create(card)
		if err != nil {
			return response.WrapError(err, http.StatusInternalServerError, "failed to create cards")
		}

		tags := r.Form["tags"]

		for _, tag := range tags {
			id, err := ub.ParseID(tag)
			if err != nil {
				msg := fmt.Sprintf("invalid tag id %s", tag)
				return response.WrapError(err, http.StatusBadRequest, msg)
			}

			ct := srs.CardTag{
				CardID: card.MetaID,
				TagID:  id,
			}

			err = conn.Create(&ct)
			if err != nil {
				return response.WrapError(err, http.StatusInternalServerError, "failed to create card tag")
			}
		}

		path, err := ub.Path("SHOW", card)
		if err != nil {
			return response.WrapError(err, http.StatusInternalServerError, "failed to build card path")
		}

		return response.Redirect{Path: path, Code: http.StatusFound}
	}
}

func Show(conn srs.Database, ub web.URLBuilder, hash string) response.Handler {
	return func(w http.ResponseWriter, r *http.Request) response.Responder {

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

func Update(conn srs.Database, ub web.URLBuilder, hash string) response.Handler {
	return func(w http.ResponseWriter, r *http.Request) response.Responder {

		return response.NewError(http.StatusInternalServerError, "not implemented")
	}
}
