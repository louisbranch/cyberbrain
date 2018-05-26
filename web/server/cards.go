package server

import (
	"net/http"

	"github.com/luizbranco/srs/web"
)

func (srv *Server) cards(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		path := r.URL.Path[len("/cards/"):]
		if path == "" {
			http.Redirect(w, r, "/decks", http.StatusFound)
			return
		}

		srv.deckShow(path, w, r)
	case "POST":
		slug := r.FormValue("deck")

		if slug == "" {
			srv.renderNotFound(w)
			return
		}

		deck := &web.Deck{}

		err := srv.Database.Get(slug, deck)
		if err != nil {
			srv.renderError(w, err)
			return
		}

		card := web.Card{
			DeckID:   deck.ID,
			ImageURL: r.FormValue("image_url"),
			AudioURL: r.FormValue("audio_url"),
			Field1:   r.FormValue("field_1"),
			Field2:   r.FormValue("field_2"),
			Field3:   r.FormValue("field_3"),
		}

		err = srv.Database.Create(&card)
		if err != nil {
			srv.renderError(w, err)
			return
		}

		http.Redirect(w, r, "/decks", http.StatusFound)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (srv *Server) newCard(w http.ResponseWriter, r *http.Request) {
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

	deck := &web.Deck{}

	err := srv.Database.Get(slug, deck)
	if err != nil {
		srv.renderError(w, err)
		return
	}

	page := web.Page{
		Title:    "New Card",
		Partials: []string{"new_card"},
		Content:  deck,
	}

	srv.render(w, page)
}
