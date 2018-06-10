package tags

import (
	"context"
	"net/http"

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

		content, err := html.RenderDeck(ub, deck, nil, nil)
		if err != nil {
			return response.WrapError(err, http.StatusInternalServerError, "failed to render deck")
		}

		page := web.Page{
			Title:    "New Tag",
			Partials: []string{"new_tag"},
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

		tag, err := html.NewTagFromForm(deck, r.Form)
		if err != nil {
			return response.WrapError(err, http.StatusBadRequest, "invalid tag form")
		}

		err = conn.Create(tag)
		if err != nil {
			return response.WrapError(err, http.StatusInternalServerError, "failed to create tag")
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

		deck, _ := middlewares.CurrentDeck(ctx)

		tag, cards, err := finder.Tag(conn, ub, hash, finder.WithCards)
		if err != nil {
			return err.(response.Error)
		}

		content, err := html.RenderTag(ub, deck, *tag, cards, true)
		if err != nil {
			return response.WrapError(err, http.StatusInternalServerError, "failed to render tag")
		}

		page := web.Page{
			Title:    "Tag",
			Partials: []string{"tag"},
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

		tag, _, err := finder.Tag(conn, ub, hash, finder.NoOption)
		if err != nil {
			return err.(response.Error)
		}

		deck, err := db.FindDeck(conn, tag.DeckID)
		if err != nil {
			return response.WrapError(err, http.StatusInternalServerError, "invalid tag deck")
		}

		newTag, err := html.NewTagFromForm(*deck, r.Form)
		if err != nil {
			return response.WrapError(err, http.StatusBadRequest, "invalid tag form")
		}

		tag.Name = newTag.Name

		err = conn.Update(tag)
		if err != nil {
			return response.WrapError(err, http.StatusBadRequest, "failed to update tag")
		}

		path, err := ub.Path("SHOW", deck)
		if err != nil {
			return response.WrapError(err, http.StatusInternalServerError, "failed to generate deck path")
		}

		return response.Redirect{Path: path, Code: http.StatusFound}
	}
}
