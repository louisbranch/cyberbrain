package home

import (
	"context"
	"net/http"

	"gitlab.com/luizbranco/srs/web/server/middlewares"
	"gitlab.com/luizbranco/srs/web/server/response"
)

func NewServeMux(renderer *middlewares.Renderer) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		renderer.Render(Index(), w, r)
	})

	return mux
}

func Index() response.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) response.Responder {
		if r.Method != "GET" || r.URL.Path != "/" {
			return response.NewError(http.StatusNotFound, r.URL.Path+" not found")
		}

		return response.Redirect{Path: "/decks/", Code: http.StatusFound}
	}
}
