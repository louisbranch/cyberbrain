package rounds

import (
	"net/http"

	"gitlab.com/luizbranco/srs/primitives"
	"gitlab.com/luizbranco/srs/web"
	"gitlab.com/luizbranco/srs/web/server/middlewares"
	"gitlab.com/luizbranco/srs/web/server/response"
)

func NewServeMux(renderer *middlewares.Renderer, db primitives.Database,
	ub web.URLBuilder, gen primitives.PracticeGenerator) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path[len("/"):]
		method := r.Method

		var handler response.Handler

		switch {
		case method == "GET" && path == "":
			handler = Index()
		case method == "GET" && path == "new":
			handler = New(db, ub, gen)
		case method == "GET":
			handler = Show(db, ub, path)
		case method == "POST" && path == "":
			handler = Create(db, ub, gen)
		case method == "POST":
			handler = Update(db, ub, path)
		}

		handler = middlewares.Authenticate(handler)

		renderer.Render(handler, w, r)
	})

	return mux
}
