package server

import (
	"net/http"

	"github.com/luizbranco/srs/web"
	"github.com/luizbranco/srs/web/models"
)

func (srv *Server) tags(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		http.Redirect(w, r, "/decks", http.StatusFound)
	case "POST":
		if err := r.ParseForm(); err != nil {
			srv.renderError(w, err)
			return
		}

		slug := r.Form.Get("deck")

		if slug == "" {
			srv.renderNotFound(w)
			return
		}

		deck, err := models.FindDeckBySlug(srv.Database, slug)
		if err != nil {
			srv.renderError(w, err)
			return
		}

		tag, err := models.NewTagFromForm(deck.MetaID, r.Form)
		if err != nil {
			srv.renderError(w, err)
			return
		}

		err = srv.Database.Create(tag)
		if err != nil {
			srv.renderError(w, err)
			return
		}

		http.Redirect(w, r, "/decks/"+slug, http.StatusFound)
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
	slug := query.Get("deck")

	if slug == "" {
		srv.renderNotFound(w)
		return
	}

	deck, err := models.FindDeckBySlug(srv.Database, slug)
	if err != nil {
		srv.renderError(w, err)
		return
	}

	page := web.Page{
		Title:    "New Tag",
		Partials: []string{"new_tag"},
		Content:  deck,
	}

	srv.render(w, page)
}
