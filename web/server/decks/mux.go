package decks

import (
	"net/http"
	"strings"

	"gitlab.com/luizbranco/srs/primitives"
	"gitlab.com/luizbranco/srs/web"
	"gitlab.com/luizbranco/srs/web/server/cards"
	"gitlab.com/luizbranco/srs/web/server/middlewares"
	"gitlab.com/luizbranco/srs/web/server/response"
)

func NewServeMux(renderer *middlewares.Renderer, db primitives.Database,
	ub web.URLBuilder) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		method := r.Method

		var handler response.Handler

		paths := strings.SplitN(r.URL.Path, "/", 4)

		deckID := paths[1]

		if len(paths) <= 2 {
			switch {
			case method == "GET" && deckID == "":
				handler = Index(db, ub)
			case method == "GET" && deckID == "new":
				handler = New(db, ub)
			case method == "GET":
				handler = Show(db, ub, deckID)
			case method == "POST" && deckID == "":
				handler = Create(db, ub)
			}

			handler = middlewares.Authenticate(handler)

			renderer.Render(handler, w, r)
			return
		}

		path := ""
		if len(paths) > 3 {
			path = paths[3]
		}

		switch paths[2] {
		case "cards":
			switch {
			case method == "GET" && path == "":
				handler = cards.Index()
			case method == "GET" && path == "new":
				handler = cards.New(db, ub)
			case method == "GET":
				handler = cards.Show(db, ub, path)
			case method == "POST" && path == "":
				handler = cards.Create(db, ub)
			case method == "POST":
				handler = cards.Update(db, ub, path)
			}
		}

		handler = middlewares.Deck(handler, db, ub, deckID)

		handler = middlewares.Authenticate(handler)

		renderer.Render(handler, w, r)
	})

	return mux
}
