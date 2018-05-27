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

		deck := &web.Deck{}

		where := web.Where{"slug": slug}

		err := srv.Database.Get(where, deck)
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

		p := web.Practice{
			DeckID: deck.ID,
			Rounds: n,
			State:  web.PracticeStateInProgress,
		}

		err = srv.Database.Create(&p)
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

	where := web.Where{"slug": slug}

	deck := &web.Deck{}

	err := srv.Database.Get(where, deck)
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

	where := web.Where{"id": id}

	p := &web.Practice{}

	err = srv.Database.Get(where, p)
	if err != nil {
		srv.renderError(w, err)
		return
	}

	where = web.Where{"id": p.DeckID}

	deck := &web.Deck{}

	err = srv.Database.Get(where, deck)
	if err != nil {
		srv.renderError(w, err)
		return
	}

	content := struct {
		Practice *web.Practice
		Deck     *web.Deck
	}{
		Practice: p,
		Deck:     deck,
	}

	page := web.Page{
		Title:    "Practice #" + slug,
		Partials: []string{"practice"},
		Content:  content,
	}

	srv.render(w, page)
}
