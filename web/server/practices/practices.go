package practices

import (
	"net/http"

	"github.com/luizbranco/srs"
	"github.com/luizbranco/srs/db"
	"github.com/luizbranco/srs/web"
	"github.com/luizbranco/srs/web/html"
	"github.com/luizbranco/srs/web/server/response"
)

func Index() response.Handler {
	return func(w http.ResponseWriter, r *http.Request) response.Responder {
		return response.Redirect{Path: "decks", Code: http.StatusFound}
	}
}

func New(conn srs.Database, ub web.URLBuilder) response.Handler {
	return func(w http.ResponseWriter, r *http.Request) response.Responder {

		query := r.URL.Query()
		hash := query.Get("deck")

		deck, err := findDeck(conn, ub, hash)
		if err != nil {
			return err.(response.Error)
		}

		content, err := html.RenderDeck(*deck, ub)
		if err != nil {
			return response.NewError(err, http.StatusInternalServerError, "failed to render deck")
		}

		page := web.Page{
			Title:    "New Practice",
			Partials: []string{"new_practice"},
			Content:  content,
		}

		return response.NewContent(page)
	}
}

func Create(conn srs.Database, ub web.URLBuilder) response.Handler {
	return func(w http.ResponseWriter, r *http.Request) response.Responder {
		if err := r.ParseForm(); err != nil {
			return response.NewError(err, http.StatusBadRequest, "invalid form")
		}

		hash := r.Form.Get("deck")

		deck, err := findDeck(conn, ub, hash)
		if err != nil {
			return err.(response.Error)
		}

		p, err := html.NewPracticeFromForm(*deck, r.Form, ub)
		if err != nil {
			return response.NewError(err, http.StatusBadRequest, "invalid practice values")
		}

		err = conn.Create(p)
		if err != nil {
			return response.NewError(err, http.StatusInternalServerError, "failed to create practice")
		}

		path, err := ub.Path("SHOW", p)
		if err != nil {
			return response.NewError(err, http.StatusInternalServerError, "failed to generate practice path")
		}

		return response.Redirect{Path: path, Code: http.StatusFound}
	}
}

func Show(conn srs.Database, ub web.URLBuilder, hash string) response.Handler {
	return func(w http.ResponseWriter, r *http.Request) response.Responder {

		id, err := ub.ParseID(hash)
		if err != nil {
			return response.NewError(err, http.StatusNotFound, "invalid practice id")
		}

		p, err := db.FindPractice(conn, id)
		if err != nil {
			return response.NewError(err, http.StatusNotFound, "wrong practice id")
		}

		deck, err := db.FindDeck(conn, p.DeckID)
		if err != nil {
			return response.NewError(err, http.StatusInternalServerError, "invalid deck id")
		}

		p.Deck = deck

		content, err := html.RenderPractice(*p, ub)
		if err != nil {
			return response.NewError(err, http.StatusInternalServerError, "failed to render practice")
		}

		page := web.Page{
			Title:    "Practice",
			Partials: []string{"practice"},
			Content:  content,
		}

		return response.NewContent(page)
	}
}

func findDeck(conn srs.Database, ub web.URLBuilder, hash string) (*srs.Deck, error) {

	id, err := ub.ParseID(hash)
	if err != nil {
		return nil, response.NewError(err, http.StatusBadRequest, "invalid deck id")
	}

	deck, err := db.FindDeck(conn, id)
	if err != nil {
		return nil, response.NewError(err, http.StatusBadRequest, "wrong deck id")
	}

	tags, err := db.FindTags(conn, id)
	if err != nil {
		return nil, response.NewError(err, http.StatusInternalServerError, "failed to find deck tags")
	}

	deck.Tags = tags

	return deck, nil
}
