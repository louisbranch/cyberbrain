package decks

import (
	"net/http"

	"gitlab.com/luizbranco/srs/db"
	"gitlab.com/luizbranco/srs/primitives"
	"gitlab.com/luizbranco/srs/web"
	"gitlab.com/luizbranco/srs/web/html"
	"gitlab.com/luizbranco/srs/web/server/finder"
	"gitlab.com/luizbranco/srs/web/server/response"
)

func Index(conn primitives.Database, ub web.URLBuilder) response.Handler {
	return func(w http.ResponseWriter, r *http.Request, user *primitives.User) response.Responder {

		if user != nil {
			return response.Redirect{Path: "/login", Code: http.StatusFound}
		}

		decks, err := db.FindDecks(conn, user.ID())
		if err != nil {
			return response.WrapError(err, http.StatusInternalServerError, "failed to find decks")
		}

		var content []*html.Deck

		for _, d := range decks {
			dr, err := html.RenderDeck(d, ub)
			if err != nil {
				return response.WrapError(err, http.StatusInternalServerError, "failed to render deck")
			}

			content = append(content, dr)
		}

		page := web.Page{
			Title:      "Decks",
			ActiveMenu: "decks",
			Partials:   []string{"decks"},
			Content:    content,
		}

		return response.NewContent(page)
	}
}

func New(conn primitives.Database, ub web.URLBuilder) response.Handler {
	return func(w http.ResponseWriter, r *http.Request, user *primitives.User) response.Responder {

		page := web.Page{
			Title:      "New Deck",
			ActiveMenu: "decks",
			Partials:   []string{"new_deck"},
		}

		return response.NewContent(page)
	}
}

func Create(conn primitives.Database, ub web.URLBuilder) response.Handler {
	return func(w http.ResponseWriter, r *http.Request, user *primitives.User) response.Responder {

		if err := r.ParseForm(); err != nil {
			return response.WrapError(err, http.StatusBadRequest, "invalid form")
		}

		deck, err := html.NewDeckFromForm(r.Form)
		if err != nil {
			return response.WrapError(err, http.StatusBadRequest, "invalid deck form")
		}

		err = conn.Create(deck)
		if err != nil {
			return response.WrapError(err, http.StatusInternalServerError, "failed to create deck")
		}

		path, err := ub.Path("SHOW", deck)
		if err != nil {
			return response.WrapError(err, http.StatusInternalServerError, "failed to generate deck path")
		}

		return response.Redirect{Path: path, Code: http.StatusFound}
	}
}

func Show(conn primitives.Database, ub web.URLBuilder, hash string) response.Handler {
	return func(w http.ResponseWriter, r *http.Request, user *primitives.User) response.Responder {

		deck, err := finder.Deck(conn, ub, hash, finder.WithTags|finder.WithCards)
		if err != nil {
			return err.(response.Error)
		}

		content, err := html.RenderDeck(*deck, ub)
		if err != nil {
			return response.WrapError(err, http.StatusInternalServerError, "failed to render deck")
		}

		page := web.Page{
			Title:      deck.Name + " Deck",
			ActiveMenu: "decks",
			Partials:   []string{"deck"},
			Content:    content,
		}

		return response.NewContent(page)
	}
}
