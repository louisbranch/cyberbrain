package sessions

import (
	"context"
	"net/http"

	"gitlab.com/luizbranco/cyberbrain/db"
	"gitlab.com/luizbranco/cyberbrain/primitives"
	"gitlab.com/luizbranco/cyberbrain/web"
	"gitlab.com/luizbranco/cyberbrain/web/server/middlewares"
	"gitlab.com/luizbranco/cyberbrain/web/server/response"
)

func New(conn primitives.Database, ub web.URLBuilder) response.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) response.Responder {

		if _, ok := middlewares.CurrentUser(ctx); ok {
			return response.Redirect{Path: "/", Code: http.StatusFound}
		}

		page := web.Page{
			Title:      "Log In",
			ActiveMenu: "login",
			Partials:   []string{"new_session"},
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

		email := r.Form.Get("email")
		password := r.Form.Get("password")

		user, err := db.FindUserByEmail(conn, email)
		if err != nil {
			return response.WrapError(err, http.StatusBadRequest, "user and password combination doesn't match")
		}

		ok, err := auth.Verify(user.PasswordHash, password)
		if err != nil {
			return response.WrapError(err, http.StatusInternalServerError, "user and password combination doesn't match")
		}

		if !ok {
			return response.WrapError(err, http.StatusBadRequest, "user and password combination doesn't match")
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

func Destroy(session web.SessionManager) response.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) response.Responder {
		session.LogOut(w)
		return response.Redirect{Path: "/", Code: http.StatusFound}
	}
}
