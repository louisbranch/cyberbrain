package blitline

import (
	"net/http"
	"strings"

	"gitlab.com/luizbranco/cyberbrain/primitives"
	"gitlab.com/luizbranco/cyberbrain/web"
	"gitlab.com/luizbranco/cyberbrain/web/server/middlewares"
	"gitlab.com/luizbranco/cyberbrain/web/server/response"
)

func NewServeMux(renderer *middlewares.Renderer, db primitives.Database,
	ub web.URLBuilder) *http.ServeMux {

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		method := r.Method

		if method != "POST" {
			renderer.Render(nil, w, r)
		}

		var handler response.Handler

		paths := strings.SplitN(r.URL.Path, "/", 3)

		if len(paths) != 3 {
			renderer.Render(nil, w, r)
		}

		switch {
		case paths[1] == "cards":
			handler = PatchCard(db, ub, paths[2])
		}

		renderer.Render(handler, w, r)
	})

	return mux
}
