package server

import (
	"net/http"

	"gitlab.com/luizbranco/srs/db"
	"gitlab.com/luizbranco/srs/web"
	"gitlab.com/luizbranco/srs/web/html"
)

func (srv *Server) tags(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (srv *Server) newTag(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query()

	id, err := srv.URLBuilder.ParseID(query.Get("deck"))
	if err != nil {
		// FIXME bad request
		srv.renderNotFound(w)
		return
	}

	deck, err := db.FindDeck(srv.Database, id)
	if err != nil {
		srv.renderError(w, err)
		return
	}

	content, err := html.RenderDeck(*deck, srv.URLBuilder)
	if err != nil {
		srv.renderError(w, err)
		return
	}

	page := web.Page{
		Title:    "New Tag",
		Partials: []string{"new_tag"},
		Content:  content,
	}

	srv.render(w, page)
}
