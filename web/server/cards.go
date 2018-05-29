package server

import (
	"net/http"

	"github.com/luizbranco/srs/web"
	"github.com/luizbranco/srs/web/models"
)

func (srv *Server) cards(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		slug := r.URL.Path[len("/cards/"):]
		if slug == "" {
			http.Redirect(w, r, "/decks", http.StatusFound)
			return
		}

		srv.cardShow(slug, w, r)
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
			srv.renderNotFound(w)
			return
		}

		card, err := models.NewCardFromForm(deck.ID, r.Form)
		if err != nil {
			srv.renderError(w, err)
			return
		}

		err = srv.Database.Create(card)
		if err != nil {
			srv.renderError(w, err)
			return
		}

		tags := r.Form["tags"]

		for _, tag := range tags {
			ct := models.CardTag{
				CardID: card.ID,
				TagID:  web.ID(tag),
			}

			err = srv.Database.Create(&ct)
			if err != nil {
				srv.renderError(w, err)
				return
			}
		}

		http.Redirect(w, r, "/decks/"+slug, http.StatusFound)
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

	deck, err := models.FindDeckBySlug(srv.Database, slug)
	if err != nil {
		// FIXME
		srv.renderError(w, err)
		return
	}

	tags, err := models.FindTagsByDeckID(srv.Database, deck.ID)
	if err != nil {
		// FIXME
		srv.renderError(w, err)
		return
	}

	deck.Tags = tags

	page := web.Page{
		Title:    "New Card",
		Partials: []string{"new_card"},
		Content:  deck,
	}

	srv.render(w, page)
}

func (srv *Server) cardShow(slug string, w http.ResponseWriter, r *http.Request) {
	card, err := models.FindCardBySlug(srv.Database, slug)
	if err != nil {
		srv.renderError(w, err)
		return
	}

	deck, err := models.FindDeckByID(srv.Database, card.DeckID)
	if err != nil {
		srv.renderError(w, err)
		return
	}

	tags, err := models.FindTagsByCardID(srv.Database, card.ID)
	if err != nil {
		srv.renderError(w, err)
		return
	}

	card.Tags = tags

	content := struct {
		Card *models.Card
		Deck *models.Deck
	}{
		Card: card,
		Deck: deck,
	}

	page := web.Page{
		Title:    "Card #" + slug,
		Partials: []string{"card"},
		Content:  content,
	}

	srv.render(w, page)
}
