package decks

import (
	"context"
	"net/http"

	"gitlab.com/luizbranco/cyberbrain/db"
	"gitlab.com/luizbranco/cyberbrain/primitives"
	"gitlab.com/luizbranco/cyberbrain/web"
	"gitlab.com/luizbranco/cyberbrain/web/html"
	"gitlab.com/luizbranco/cyberbrain/web/server/finder"
	"gitlab.com/luizbranco/cyberbrain/web/server/middlewares"
	"gitlab.com/luizbranco/cyberbrain/web/server/response"
)

func Index(conn primitives.Database, ub web.URLBuilder) response.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) response.Responder {

		user, _ := middlewares.CurrentUser(ctx)

		decks, err := db.FindDecks(conn, user.ID())
		if err != nil {
			return response.WrapError(err, http.StatusInternalServerError, "failed to find decks")
		}

		var content []*html.Deck

		for _, d := range decks {
			dr, err := html.RenderDeck(ub, d, nil, nil)
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
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) response.Responder {

		page := web.Page{
			Title:      "New Deck",
			ActiveMenu: "decks",
			Partials:   []string{"new_deck"},
		}

		return response.NewContent(page)
	}
}

func Create(conn primitives.Database, ub web.URLBuilder) response.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) response.Responder {

		if err := r.ParseForm(); err != nil {
			return response.WrapError(err, http.StatusBadRequest, "invalid form")
		}

		deck, err := html.NewDeckFromForm(r.Form)
		if err != nil {
			return response.WrapError(err, http.StatusBadRequest, "invalid deck form")
		}

		user, _ := middlewares.CurrentUser(ctx)

		deck.UserID = user.ID()

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
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) response.Responder {

		deck, cards, tags, err := finder.Deck(conn, ub, hash, finder.WithTags|finder.WithCards)
		if err != nil {
			return err.(response.Error)
		}

		content, err := html.RenderDeck(ub, *deck, cards, tags)
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

func Edit(conn primitives.Database, ub web.URLBuilder) response.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) response.Responder {

		deck := middlewares.CurrentDeck(ctx)

		content, err := html.RenderDeck(ub, deck, nil, nil)
		if err != nil {
			return response.WrapError(err, http.StatusInternalServerError, "failed to render deck")
		}

		page := web.Page{
			Title:      deck.Name + " Deck",
			ActiveMenu: "decks",
			Partials:   []string{"edit_deck"},
			Content:    content,
		}

		return response.NewContent(page)
	}
}

func Update(conn primitives.Database, ub web.URLBuilder, hash string) response.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) response.Responder {

		if err := r.ParseForm(); err != nil {
			return response.WrapError(err, http.StatusBadRequest, "invalid form")
		}

		deck, _, _, err := finder.Deck(conn, ub, hash, finder.NoOption)
		if err != nil {
			return err.(response.Error)
		}

		deck.Name = r.Form.Get("name")
		deck.Description = r.Form.Get("description")

		if deck.Name == "" {
			return response.NewError(http.StatusBadRequest, "deck name cannot be empty")
		}

		err = conn.Update(deck)
		if err != nil {
			return response.WrapError(err, http.StatusBadRequest, "failed to update deck")
		}

		path, err := ub.Path("SHOW", deck)
		if err != nil {
			return response.WrapError(err, http.StatusInternalServerError, "failed to generate deck path")
		}

		return response.Redirect{Path: path, Code: http.StatusFound}
	}
}
