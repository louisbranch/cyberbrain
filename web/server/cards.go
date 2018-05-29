package server

import (
	"net/http"
	"strconv"

	"github.com/luizbranco/srs/web"
	"github.com/pkg/errors"
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
		err := r.ParseForm()
		if err != nil {
			srv.renderError(w, err)
			return
		}

		slug := r.FormValue("deck")
		if slug == "" {
			srv.renderNotFound(w)
			return
		}

		deck, err := FindDeckBySlug(srv.Database, slug)
		if err != nil {
			srv.renderNotFound(w)
			return
		}

		card := &web.Card{
			DeckID:        deck.ID,
			ImageURL:      r.FormValue("image_url"),
			AudioURL:      r.FormValue("audio_url"),
			Definition:    r.FormValue("definition"),
			AltDefinition: r.FormValue("alt_definition"),
			Pronunciation: r.FormValue("pronunciation"),
		}

		err = srv.Database.Create(card)
		if err != nil {
			srv.renderError(w, err)
			return
		}

		tags := r.Form["tags"]

		for _, tag := range tags {
			tid, err := strconv.ParseUint(tag, 10, 64)
			if err != nil {
				err = errors.Wrapf(err, "invalid tag id %d", tag)
				srv.renderError(w, err)
				return
			}

			ct := &web.CardTag{
				CardID: card.ID,
				TagID:  uint(tid),
			}

			err = srv.Database.Create(ct)
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

	deck, err := FindDeckBySlug(srv.Database, slug)
	if err != nil {
		// FIXME
		srv.renderError(w, err)
		return
	}

	tags, err := FindTagsByDeckID(srv.Database, deck.ID)
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
	id, err := strconv.ParseUint(slug, 10, 64)
	if err != nil {
		err = errors.Wrapf(err, "invalid card id %d", slug)
		srv.renderError(w, err)
		return
	}

	card, err := FindCardByID(srv.Database, id)
	if err != nil {
		srv.renderError(w, err)
		return
	}

	deck, err := FindDeckByID(srv.Database, card.DeckID)
	if err != nil {
		srv.renderError(w, err)
		return
	}

	tags, err := FindTagsByCardID(srv.Database, id)
	if err != nil {
		srv.renderError(w, err)
		return
	}

	card.Tags = tags

	content := struct {
		Card *web.Card
		Deck *web.Deck
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
