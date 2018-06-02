package server

import (
	"net/http"
	"sort"

	"github.com/luizbranco/srs/web"
	"github.com/luizbranco/srs/web/models"
)

func (srv *Server) decks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		path := r.URL.Path[len("/decks/"):]
		if path == "" {
			srv.decksList(w, r)
			return
		}

		id, err := srv.URLBuilder.ParseID(path)
		if err != nil {
			srv.renderNotFound(w)
			return
		}

		srv.deckShow(id, w, r)
	case "POST":

		if err := r.ParseForm(); err != nil {
			srv.renderError(w, err)
			return
		}

		deck, err := models.NewDeckFromForm(r.Form)
		if err != nil {
			srv.renderError(w, err)
			return
		}

		err = srv.Database.Create(deck)
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
	decks, err := models.FindDecks(srv.Database)
	if err != nil {
		srv.renderError(w, err)
		return
	}

	var content []*models.DeckRendered

	for _, d := range decks {
		dr, err := d.Render(srv.URLBuilder)
		if err != nil {
			srv.renderError(w, err)
			return
		}

		content = append(content, dr)
	}

	page := web.Page{
		Title:      "Decks",
		ActiveMenu: "decks",
		Partials:   []string{"decks"},
		Content:    content,
	}

	srv.render(w, page)
}

func (srv *Server) deckShow(id web.ID, w http.ResponseWriter, r *http.Request) {
	deck, err := models.FindDeck(srv.Database, id)
	if err != nil {
		srv.renderError(w, err)
		return
	}

	cards, err := models.FindCardsByDeck(srv.Database, deck.MetaID)
	if err != nil {
		srv.renderError(w, err)
		return
	}

	tags, err := models.FindTagsByDeckID(srv.Database, deck.MetaID)
	if err != nil {
		srv.renderError(w, err)
		return
	}

	sort.Slice(tags, func(i, j int) bool {
		return tags[i].Name < tags[j].Name
	})

	deck.Cards = cards
	deck.Tags = tags

	content, err := deck.Render(srv.URLBuilder)
	if err != nil {
		srv.renderError(w, err)
		return
	}

	page := web.Page{
		Title:      deck.Name + " Deck",
		ActiveMenu: "decks",
		Partials:   []string{"deck"},
		Content:    content,
	}

	srv.render(w, page)
}
