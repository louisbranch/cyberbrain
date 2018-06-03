package server

import (
	"net/http"

	"github.com/luizbranco/srs"
	"github.com/luizbranco/srs/db"
	"github.com/luizbranco/srs/web"
	"github.com/luizbranco/srs/web/html"
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

		id, err := srv.URLBuilder.ParseID(path)
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

		id, err := srv.URLBuilder.ParseID(r.Form.Get("deck"))
		if err != nil {
			// FIXME bad request
			srv.renderNotFound(w)
			return
		}

		deck, err := db.FindDeck(srv.Database, id)
		if err != nil {
			srv.renderNotFound(w)
			return
		}

		card, err := html.NewCardFromForm(*deck, r.Form)
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
			id, err := srv.URLBuilder.ParseID(tag)
			if err != nil {
				err = errors.Wrapf(err, "invalid tag id %s", tag)
				srv.renderError(w, err)
				return
			}

			ct := srs.CardTag{
				CardID: card.MetaID,
				TagID:  id,
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
	id, err := srv.URLBuilder.ParseID(query.Get("deck"))
	if err != nil {
		srv.renderNotFound(w)
		return
	}

	deck, err := db.FindDeck(srv.Database, id)
	if err != nil {
		// FIXME bad request
		srv.renderError(w, err)
		return
	}

	tags, err := db.FindTags(srv.Database, deck.MetaID)
	if err != nil {
		// FIXME bad request
		srv.renderError(w, err)
		return
	}

	deck.Tags = tags

	content, err := html.RenderDeck(*deck, srv.URLBuilder)
	if err != nil {
		srv.renderError(w, err)
		return
	}

	page := web.Page{
		Title:    "New Card",
		Partials: []string{"new_card"},
		Content:  content,
	}

	srv.render(w, page)
}

func (srv *Server) cardShow(id srs.ID, w http.ResponseWriter, r *http.Request) {
	card, err := db.FindCard(srv.Database, id)
	if err != nil {
		srv.renderError(w, err)
		return
	}

	deck, err := db.FindDeck(srv.Database, card.DeckID)
	if err != nil {
		srv.renderError(w, err)
		return
	}

	tags, err := db.FindTagsByCard(srv.Database, card.ID())
	if err != nil {
		srv.renderError(w, err)
		return
	}

	card.Tags = tags

	dr, err := html.RenderDeck(*deck, srv.URLBuilder)
	if err != nil {
		srv.renderError(w, err)
		return
	}

	cr, err := html.RenderCard(*card, srv.URLBuilder)
	if err != nil {
		srv.renderError(w, err)
		return
	}

	content := struct {
		Card *html.Card
		Deck *html.Deck
	}{
		Card: cr,
		Deck: dr,
	}

	page := web.Page{
		Title:    "Card",
		Partials: []string{"card"},
		Content:  content,
	}

	srv.render(w, page)
}
