package home

import (
	"context"
	"log"
	"net/http"

	"gitlab.com/luizbranco/cyberbrain/web"
	"gitlab.com/luizbranco/cyberbrain/web/server/middlewares"
	"gitlab.com/luizbranco/cyberbrain/web/server/response"
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
			log.Printf("[404] %s", r.URL.Path)
			return response.NewError(http.StatusNotFound, r.URL.Path+" not found")
		}

		if _, ok := middlewares.CurrentUser(ctx); ok {
			return response.Redirect{Path: "/decks/", Code: http.StatusFound}
		}

		page := web.Page{
			Title:      "CyberBrain.app",
			ActiveMenu: "home",
			Partials:   []string{"home"},
		}

		return response.NewContent(page)
	}
}
