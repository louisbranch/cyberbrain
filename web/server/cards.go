package server

import (
	"net/http"
	"strconv"

	"github.com/luizbranco/srs/web"
	"github.com/luizbranco/srs/web/models"
	"github.com/pkg/errors"
)

func (srv *Server) cards(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		path := r.URL.Path[len("/cards/"):]
		if path == "" {
			http.Redirect(w, r, "/decks", http.StatusFound)
			return
		}

		id, err := srv.URLBuilder.ID(path)
		if err != nil {
			srv.renderNotFound(w)
			return
		}

		srv.cardShow(id, w, r)
	case "POST":
		if err := r.ParseForm(); err != nil {
			srv.renderError(w, err)
			return
		}

		id, err := srv.URLBuilder.ID(r.Form.Get("deck"))
		if err != nil {
			// FIXME bad request
			srv.renderNotFound(w)
			return
		}

		deck, err := models.FindDeck(srv.Database, id)
		if err != nil {
			srv.renderNotFound(w)
			return
		}

		card, err := models.NewCardFromForm(deck.MetaID, r.Form)
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
			id, err := strconv.Atoi(tag)
			if err != nil {
				err = errors.Wrapf(err, "invalid tag id %s", tag)
				srv.renderError(w, err)
				return
			}

			ct := models.CardTag{
				CardID: card.MetaID,
				TagID:  web.ID(id),
			}

			err = srv.Database.Create(&ct)
			if err != nil {
				srv.renderError(w, err)
				return
			}
		}

		path, err := srv.URLBuilder.Path("SHOW", card)
		if err != nil {
			srv.renderError(w, err)
			return
		}

		http.Redirect(w, r, path, http.StatusFound)
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
	id, err := srv.URLBuilder.ID(query.Get("deck"))
	if err != nil {
		srv.renderNotFound(w)
		return
	}

	deck, err := models.FindDeck(srv.Database, id)
	if err != nil {
		// FIXME bad request
		srv.renderError(w, err)
		return
	}

	tags, err := models.FindTagsByDeckID(srv.Database, deck.MetaID)
	if err != nil {
		// FIXME bad request
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

func (srv *Server) cardShow(id web.ID, w http.ResponseWriter, r *http.Request) {
	card, err := models.FindCard(srv.Database, id)
	if err != nil {
		srv.renderError(w, err)
		return
	}

	deck, err := models.FindDeck(srv.Database, card.DeckID)
	if err != nil {
		srv.renderError(w, err)
		return
	}

	tags, err := models.FindTagsByCard(srv.Database, card.MetaID)
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
		Title:    "Card",
		Partials: []string{"card"},
		Content:  content,
	}

	srv.render(w, page)
}
