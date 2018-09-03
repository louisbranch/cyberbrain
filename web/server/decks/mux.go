package decks

import (
	"net/http"
	"strings"

	"gitlab.com/luizbranco/cyberbrain/primitives"
	"gitlab.com/luizbranco/cyberbrain/web"
	"gitlab.com/luizbranco/cyberbrain/web/server/cards"
	"gitlab.com/luizbranco/cyberbrain/web/server/middlewares"
	"gitlab.com/luizbranco/cyberbrain/web/server/response"
	"gitlab.com/luizbranco/cyberbrain/web/server/reviews"
	"gitlab.com/luizbranco/cyberbrain/web/server/tags"
	"gitlab.com/luizbranco/cyberbrain/worker"
)

func NewServeMux(renderer *middlewares.Renderer, db primitives.Database,
	ub web.URLBuilder, resizer worker.ImageResizer) *http.ServeMux {

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		method := r.Method

		var handler response.Handler

		paths := strings.Split(r.URL.Path, "/")

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
			case method == "POST":
				handler = Update(db, ub, deckID)
			}

			if handler != nil {
				handler = middlewares.Authenticate(handler)
			}

			renderer.Render(handler, w, r)
			return
		}

		path := ""
		if len(paths) > 3 {
			path = paths[3]
		}

		switch paths[2] {
		case "edit":
			if method == "GET" && path == "" {
				handler = Edit(db, ub)
			}

		case "cards":
			switch {
			case method == "GET" && path == "":
				handler = cards.Index()
			case method == "GET" && path == "new":
				handler = cards.New(db, ub)
			case method == "GET":
				handler = cards.Show(db, ub, path)
			case method == "POST" && path == "":
				handler = cards.Create(db, ub, resizer)
			case method == "POST":
				handler = cards.Update(db, ub, path)
			}

		case "tags":
			switch {
			case method == "GET" && path == "":
				handler = tags.Index()
			case method == "GET" && path == "new":
				handler = tags.New(db, ub)
			case method == "GET":
				handler = tags.Show(db, ub, path)
			case method == "POST" && path == "":
				handler = tags.Create(db, ub)
			case method == "POST":
				handler = tags.Update(db, ub, path)
			}

		case "reviews":
			switch {
			case method == "GET" && path == "":
				handler = reviews.Index()
			case method == "GET" && path == "new":
				handler = reviews.New(db, ub)
			case method == "GET":
				handler = reviews.Show(db, ub, path)
			case method == "POST":
				handler = reviews.Create(db, ub)
			}
		}

		if handler != nil {
			handler = middlewares.Deck(handler, db, ub, deckID)

			handler = middlewares.Authenticate(handler)
		}

		renderer.Render(handler, w, r)
	})

	return mux
}
