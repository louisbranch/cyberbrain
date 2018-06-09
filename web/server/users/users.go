package users

import (
	"net/http"

	"gitlab.com/luizbranco/srs/primitives"
	"gitlab.com/luizbranco/srs/web"
	"gitlab.com/luizbranco/srs/web/html"
	"gitlab.com/luizbranco/srs/web/server/response"
)

func New(conn primitives.Database, ub web.URLBuilder) response.Handler {
	return func(w http.ResponseWriter, r *http.Request) response.Responder {

		page := web.Page{
			Title:      "Sign Up",
			ActiveMenu: "signup",
			Partials:   []string{"new_user"},
		}

		return response.NewContent(page)
	}
}

func Create(conn primitives.Database, ub web.URLBuilder, auth primitives.Authenticator) response.Handler {
	return func(w http.ResponseWriter, r *http.Request) response.Responder {

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

		path, err := ub.Path("INDEX", primitives.Deck{})
		if err != nil {
			return response.WrapError(err, http.StatusInternalServerError, "failed to generate user path")
		}

		return response.Redirect{Path: path, Code: http.StatusFound}
	}
}
