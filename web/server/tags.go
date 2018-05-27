package server

import (
	"net/http"

	"github.com/luizbranco/srs/web"
)

func (srv *Server) tags(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		http.Redirect(w, r, "/decks", http.StatusFound)
	case "POST":
		slug := r.FormValue("deck")

		if slug == "" {
			srv.renderNotFound(w)
			return
		}

		deck := &web.Deck{}

		where := web.Where{
			"slug": slug,
		}

		err := srv.Database.Get(where, deck)
		if err != nil {
			srv.renderError(w, err)
			return
		}

		tag := web.Tag{
			DeckID: deck.ID,
			Name:   r.FormValue("name"),
		}

		err = srv.Database.Create(&tag)
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

	where := web.Where{
		"slug": slug,
	}

	deck := &web.Deck{}

	err := srv.Database.Get(where, deck)
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
