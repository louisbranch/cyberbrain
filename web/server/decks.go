package server

import (
	"fmt"
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

		deck := web.Deck{
			Name:        r.FormValue("name"),
			Description: r.FormValue("description"),
			ImageURL:    r.FormValue("image_url"),
			Field1:      r.FormValue("field_1"),
			Field2:      r.FormValue("field_2"),
			Field3:      r.FormValue("field_3"),
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
		Title:      "New Deck",
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

	where := fmt.Sprintf("deck_id = %d", deck.ID)

	cards := &web.Cards{}

	err = srv.Database.Query(where, cards)
	if err != nil {
		srv.renderError(w, err)
		return
	}

	tags := &web.Tags{}

	err = srv.Database.Query(where, tags)
	if err != nil {
		srv.renderError(w, err)
		return
	}

	deck.Tags = *tags

	page := web.Page{
		Title:      deck.Name + " Deck",
		ActiveMenu: "decks",
		Partials:   []string{"deck"},
		Content:    deck,
	}

	srv.render(w, page)
}
