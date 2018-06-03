package server

import (
	"net/http"

	"github.com/luizbranco/srs"
	"github.com/luizbranco/srs/db"
	"github.com/luizbranco/srs/web"
	"github.com/luizbranco/srs/web/html"
)

func (srv *Server) practice(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		path := r.URL.Path[len("/practices/"):]
		if path == "" {
			http.Redirect(w, r, "/decks", http.StatusFound)
			return
		}

		id, err := srv.URLBuilder.ParseID(path)
		if err != nil {
			srv.renderNotFound(w)
			return
		}

		srv.practiceShow(id, w, r)
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
			srv.renderError(w, err)
			return
		}

		p, err := html.NewPracticeFromForm(deck.MetaID, r.Form)
		if err != nil {
			srv.renderError(w, err)
			return
		}

		err = srv.Database.Create(p)
		if err != nil {
			srv.renderError(w, err)
			return
		}

		path, err := srv.URLBuilder.Path("INDEX", p)
		if err != nil {
			srv.renderError(w, err)
			return
		}

		http.Redirect(w, r, path, http.StatusFound)
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

	id, err := srv.URLBuilder.ParseID(query.Get("deck"))
	if err != nil {
		// FIXME bad request
		srv.renderNotFound(w)
		return
	}

	deck, err := db.FindDeck(srv.Database, id)
	if err != nil {
		srv.renderError(w, err)
		return
	}

	content, err := html.RenderDeck(*deck, srv.URLBuilder)
	if err != nil {
		srv.renderError(w, err)
		return
	}

	page := web.Page{
		Title:    "New Practice",
		Partials: []string{"new_practice"},
		Content:  content,
	}

	srv.render(w, page)
}

func (srv *Server) practiceShow(id srs.ID, w http.ResponseWriter, r *http.Request) {
	p, err := db.FindPractice(srv.Database, id)
	if err != nil {
		srv.renderError(w, err)
		return
	}

	deck, err := db.FindDeck(srv.Database, p.DeckID)
	if err != nil {
		srv.renderError(w, err)
		return
	}

	content := struct {
		Practice *srs.Practice
		Deck     *srs.Deck
		Card     *srs.Card
		Round    *srs.PracticeRound
	}{
		Practice: p,
		Deck:     deck,
	}

	page := web.Page{
		Title:    "Practice",
		Partials: []string{"practice"},
		Content:  content,
	}

	if p.Done {
		srv.render(w, page)
		return
	}

	n, err := db.CountPracticeRounds(srv.Database, p.MetaID)
	if err != nil {
		srv.renderError(w, err)
		return
	}

	if n < p.Rounds {
		card, err := db.FindRandomCard(srv.Database, p.DeckID)
		if err != nil {
			srv.renderError(w, err)
			return
		}

		// FIXME
		pr := &srs.PracticeRound{
			PracticeID: p.ID(),
		}

		err = srv.Database.Create(pr)
		if err != nil {
			srv.renderError(w, err)
			return
		}

		content.Card = card
		content.Round = pr
	} else {
		pr, err := db.FindPracticeRound(srv.Database, p.MetaID, n)
		if err != nil {
			srv.renderError(w, err)
			return
		}

		content.Round = pr
	}

	page = web.Page{
		Title:    "Practice",
		Partials: []string{"practice"},
		Content:  content,
	}

	srv.render(w, page)
}
