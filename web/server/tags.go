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

		id, err := srv.URLBuilder.ParseID(r.Form.Get("deck"))
		if err != nil {
			// FIXME bad request
			srv.renderNotFound(w)
			return
		}

		deck, err := models.FindDeck(srv.Database, id)
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

		path, err := srv.URLBuilder.Path("SHOW", deck)
		if err != nil {
			srv.renderError(w, err)
			return
		}

		http.Redirect(w, r, path, http.StatusFound)
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

	deck, err := models.FindDeck(srv.Database, id)
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
