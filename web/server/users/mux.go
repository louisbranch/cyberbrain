package users

import (
	"net/http"

	"gitlab.com/luizbranco/cyberbrain/primitives"
	"gitlab.com/luizbranco/cyberbrain/web"
	"gitlab.com/luizbranco/cyberbrain/web/server/middlewares"
	"gitlab.com/luizbranco/cyberbrain/web/server/response"
)

func NewSignupMux(renderer *middlewares.Renderer, db primitives.Database,
	ub web.URLBuilder, auth primitives.Authenticator) *http.ServeMux {

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path[len("/"):]
		method := r.Method

		var handler response.Handler

		switch {
		case method == "GET" && path == "":
			handler = New(db, ub)
		case method == "POST" && path == "":
			handler = Create(db, ub, auth, renderer.SessionManager)
		}

		renderer.Render(handler, w, r)
	})

	return mux
}
