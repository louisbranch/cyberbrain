package server

import (
	"fmt"
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

		deck := &web.Deck{}

		where := web.Where{
			"slug": slug,
		}

		err = srv.Database.Get(where, deck)
		if err != nil {
			srv.renderError(w, err)
			return
		}

		card := &web.Card{
			DeckID:   deck.ID,
			ImageURL: r.FormValue("image_url"),
			AudioURL: r.FormValue("audio_url"),
			Field1:   r.FormValue("field_1"),
			Field2:   r.FormValue("field_2"),
			Field3:   r.FormValue("field_3"),
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

	where := web.Where{"slug": slug}

	deck := &web.Deck{}

	err := srv.Database.Get(where, deck)
	if err != nil {
		srv.renderError(w, err)
		return
	}

	where = web.Where{"deck_id": deck.ID}

	tags := &web.Tags{}

	err = srv.Database.Query(where, tags)
	if err != nil {
		srv.renderError(w, err)
		return
	}

	deck.Tags = *tags

	page := web.Page{
		Title:    "New Card",
		Partials: []string{"new_card"},
		Content:  deck,
	}

	srv.render(w, page)
}

func (srv *Server) cardShow(slug string, w http.ResponseWriter, r *http.Request) {
	where := web.Where{"slug": slug}

	card := &web.Card{}

	err := srv.Database.Get(where, card)
	if err != nil {
		srv.renderError(w, err)
		return
	}

	where = web.Where{"id": card.DeckID}

	deck := &web.Deck{}

	err = srv.Database.Get(where, deck)
	if err != nil {
		srv.renderError(w, err)
		return
	}

	id := fmt.Sprintf("%d", card.ID)

	wRaw := `SELECT t.id, t.deck_id, t.slug, name FROM tags t
	LEFT JOIN card_tags ct ON t.id = ct.tag_id
	WHERE ct.card_id = ` + id + ";"

	tags := &web.Tags{}
	err = srv.Database.QueryRaw(wRaw, tags)
	if err != nil {
		srv.renderError(w, err)
		return
	}

	card.Tags = *tags

	content := struct {
		Card *web.Card
		Deck *web.Deck
	}{
		Card: card,
		Deck: deck,
	}

	page := web.Page{
		Title:    "Card #" + card.Slug,
		Partials: []string{"card"},
		Content:  content,
	}

	srv.render(w, page)
}
