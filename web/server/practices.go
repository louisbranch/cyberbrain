package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/luizbranco/srs/web"
	"github.com/pkg/errors"
)

func (srv *Server) practice(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		slug := r.URL.Path[len("/practices/"):]
		if slug == "" {
			http.Redirect(w, r, "/decks", http.StatusFound)
			return
		}

		srv.practiceShow(slug, w, r)
	case "POST":
		slug := r.FormValue("deck")

		if slug == "" {
			srv.renderNotFound(w)
			return
		}

		deck, err := FindDeckBySlug(srv.Database, slug)
		if err != nil {
			srv.renderError(w, err)
			return
		}

		rounds := r.FormValue("rounds")
		n, err := strconv.Atoi(rounds)
		if err != nil {
			err = errors.Wrap(err, "invalid number of rounds")
			srv.renderError(w, err)
			return
		}

		p := &web.Practice{
			DeckID: deck.ID,
			Rounds: n,
			State:  web.PracticeStateInProgress,
		}

		err = srv.Database.Create(p)
		if err != nil {
			srv.renderError(w, err)
			return
		}

		id := fmt.Sprintf("%d", p.ID)

		http.Redirect(w, r, "/practices/"+id, http.StatusFound)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (srv *Server) newPractice(w http.ResponseWriter, r *http.Request) {
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
		srv.renderError(w, err)
		return
	}

	page := web.Page{
		Title:    "New Practice",
		Partials: []string{"new_practice"},
		Content:  deck,
	}

	srv.render(w, page)
}

func (srv *Server) practiceShow(slug string, w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(slug, 10, 64)
	if err != nil {
		err = errors.Wrapf(err, "invalid practice id %d", slug)
		srv.renderError(w, err)
		return
	}

	p, err := FindPracticeByID(srv.Database, id)
	if err != nil {
		srv.renderError(w, err)
		return
	}

	deck, err := FindDeckByID(srv.Database, p.DeckID)
	if err != nil {
		srv.renderError(w, err)
		return
	}

	content := struct {
		Practice *web.Practice
		Deck     *web.Deck
		Card     *web.Card
		Round    *web.PracticeRound
	}{
		Practice: p,
		Deck:     deck,
	}

	page := web.Page{
		Title:    "Practice #" + slug,
		Partials: []string{"practice"},
		Content:  content,
	}

	if p.Finished() {
		srv.render(w, page)
		return
	}

	n, err := CountPracticeRounds(srv.Database, p.ID)
	if err != nil {
		srv.renderError(w, err)
		return
	}

	if n < p.Rounds {
		card, err := FindRandomCard(srv.Database, p.DeckID)
		if err != nil {
			srv.renderError(w, err)
			return
		}

		pr := &web.PracticeRound{
			PracticeID: p.ID,
			CardID:     card.ID,
			Expect:     card.Definition,
			Round:      n + 1,
		}

		err = srv.Database.Create(pr)
		if err != nil {
			srv.renderError(w, err)
			return
		}

		content.Card = card
		content.Round = pr
	} else {
		pr, err := FindPracticeRound(srv.Database, p.ID, n)
		if err != nil {
			srv.renderError(w, err)
			return
		}

		card, err := FindCardByID(srv.Database, uint64(pr.CardID))
		if err != nil {
			srv.renderError(w, err)
			return
		}

		content.Card = card
		content.Round = pr
	}

	page = web.Page{
		Title:    "Practice #" + slug,
		Partials: []string{"practice"},
		Content:  content,
	}

	srv.render(w, page)
}
