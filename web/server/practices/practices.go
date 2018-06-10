package practices

import (
	"context"
	"net/http"
	"sort"

	"gitlab.com/luizbranco/srs/db"
	"gitlab.com/luizbranco/srs/primitives"
	"gitlab.com/luizbranco/srs/web"
	"gitlab.com/luizbranco/srs/web/html"
	"gitlab.com/luizbranco/srs/web/server/middlewares"
	"gitlab.com/luizbranco/srs/web/server/response"
)

func Index() response.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) response.Responder {
		return response.Redirect{Path: "decks", Code: http.StatusFound}
	}
}

func New(conn primitives.Database, ub web.URLBuilder) response.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) response.Responder {

		deck := middlewares.CurrentDeck(ctx)

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
			Title:    "New Practice",
			Partials: []string{"new_practice"},
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

		deck := middlewares.CurrentDeck(ctx)

		tags, err := db.FindTags(conn, deck.ID())
		if err != nil {
			return response.WrapError(err, http.StatusInternalServerError, "failed to find deck tags")
		}

		p, err := html.NewPracticeFromForm(deck, tags, r.Form, ub)
		if err != nil {
			return response.WrapError(err, http.StatusBadRequest, "invalid practice values")
		}

		err = conn.Create(p)
		if err != nil {
			return response.WrapError(err, http.StatusInternalServerError, "failed to create practice")
		}

		path, err := ub.Path("SHOW", p, deck)
		if err != nil {
			return response.WrapError(err, http.StatusInternalServerError, "failed to generate practice path")
		}

		return response.Redirect{Path: path, Code: http.StatusFound}
	}
}

func Show(conn primitives.Database, ub web.URLBuilder, hash string) response.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) response.Responder {

		id, err := ub.ParseID(hash)
		if err != nil {
			return response.WrapError(err, http.StatusNotFound, "invalid practice id")
		}

		p, err := db.FindPractice(conn, id)
		if err != nil {
			return response.WrapError(err, http.StatusNotFound, "wrong practice id")
		}

		deck := middlewares.CurrentDeck(ctx)

		content, err := html.RenderPractice(ub, deck, *p, true)
		if err != nil {
			return response.WrapError(err, http.StatusInternalServerError, "failed to render practice")
		}

		page := web.Page{
			Title:    "Practice",
			Partials: []string{"practice"},
			Content:  content,
		}

		return response.NewContent(page)
	}
}
