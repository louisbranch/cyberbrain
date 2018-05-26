package server

import (
	"net/http"

	"github.com/luizbranco/srs/web"
)

func (srv *Server) decks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		path := r.URL.Path[len("/decks/"):]
		if path == "" {
			srv.decksList(w, r)
			return
		}

		srv.deckShow(path, w, r)
	case "POST":
		name := r.FormValue("name")
		desc := r.FormValue("description")

		deck := web.Deck{
			Name:        name,
			Description: desc,
		}

		err := srv.Database.Create(&deck)
		if err != nil {
			srv.renderError(w, err)
			return
		}

		http.Redirect(w, r, "/decks", http.StatusFound)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (srv *Server) newDeck(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	page := web.Page{
		Title:      "New Decks",
		ActiveMenu: "decks",
		Partials:   []string{"new_deck"},
	}

	srv.render(w, page)
}

func (srv *Server) decksList(w http.ResponseWriter, r *http.Request) {
	decks := &web.Decks{}

	err := srv.Database.Query("", decks)
	if err != nil {
		srv.renderError(w, err)
		return
	}

	content := struct {
		Decks *web.Decks
	}{
		Decks: decks,
	}

	page := web.Page{
		Title:      "Decks",
		ActiveMenu: "decks",
		Partials:   []string{"decks"},
		Content:    content,
	}

	srv.render(w, page)
}

func (srv *Server) deckShow(slug string, w http.ResponseWriter, r *http.Request) {
	deck := &web.Deck{}

	err := srv.Database.Get(slug, deck)
	if err != nil {
		srv.renderError(w, err)
		return
	}

	page := web.Page{
		Title:      deck.Name + " Deck",
		ActiveMenu: "decks",
		Partials:   []string{"deck"},
		Content:    deck,
	}

	srv.render(w, page)
}
