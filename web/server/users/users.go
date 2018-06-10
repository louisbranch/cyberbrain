package users

import (
	"context"
	"net/http"

	"gitlab.com/luizbranco/srs/primitives"
	"gitlab.com/luizbranco/srs/web"
	"gitlab.com/luizbranco/srs/web/html"
	"gitlab.com/luizbranco/srs/web/server/middlewares"
	"gitlab.com/luizbranco/srs/web/server/response"
)

func New(conn primitives.Database, ub web.URLBuilder) response.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) response.Responder {

		if _, ok := middlewares.CurrentUser(ctx); ok {
			return response.Redirect{Path: "/", Code: http.StatusFound}
		}

		page := web.Page{
			Title:      "Sign Up",
			ActiveMenu: "signup",
			Partials:   []string{"new_user"},
		}

		return response.NewContent(page)
	}
}

func Create(conn primitives.Database, ub web.URLBuilder,
	auth primitives.Authenticator, session web.SessionManager) response.Handler {

	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) response.Responder {

		if _, ok := middlewares.CurrentUser(ctx); ok {
			return response.Redirect{Path: "/", Code: http.StatusFound}
		}

		if err := r.ParseForm(); err != nil {
			return response.WrapError(err, http.StatusBadRequest, "invalid form")
		}

		user, err := html.NewUserFromForm(r.Form, auth)
		if err != nil {
			return response.WrapError(err, http.StatusBadRequest, "invalid user form")
		}

		err = conn.Create(user)
		if err != nil {
			return response.WrapError(err, http.StatusInternalServerError, "failed to create user")
		}

		err = session.LogIn(*user, w)
		if err != nil {
			return response.WrapError(err, http.StatusInternalServerError, "failed to log in user")
		}

		path, err := ub.Path("INDEX", primitives.Deck{})
		if err != nil {
			return response.WrapError(err, http.StatusInternalServerError, "failed to generate decks path")
		}

		return response.Redirect{Path: path, Code: http.StatusFound}
	}
}
