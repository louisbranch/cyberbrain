package server

import (
	"net/http"

	"github.com/luizbranco/srs/web"
	"github.com/luizbranco/srs/web/models"
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

		deck, err := models.FindDeck(srv.Database, id)
		if err != nil {
			srv.renderError(w, err)
			return
		}

		p, err := models.NewPracticeFromForm(deck.MetaID, r.Form)
		if err != nil {
			srv.renderError(w, err)
			return
		}

		err = srv.Database.Create(p)
		if err != nil {
			srv.renderError(w, err)
			return
		}

		http.Redirect(w, r, "/practices/", http.StatusFound) // FIXME
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

	deck, err := models.FindDeck(srv.Database, id)
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

func (srv *Server) practiceShow(id web.ID, w http.ResponseWriter, r *http.Request) {
	p, err := models.FindPractice(srv.Database, id)
	if err != nil {
		srv.renderError(w, err)
		return
	}

	deck, err := models.FindDeck(srv.Database, p.DeckID)
	if err != nil {
		srv.renderError(w, err)
		return
	}

	content := struct {
		Practice *models.Practice
		Deck     *models.Deck
		Card     *models.Card
		Round    *models.PracticeRound
	}{
		Practice: p,
		Deck:     deck,
	}

	page := web.Page{
		Title:    "Practice",
		Partials: []string{"practice"},
		Content:  content,
	}

	if p.Finished() {
		srv.render(w, page)
		return
	}

	n, err := models.CountPracticeRounds(srv.Database, p.MetaID)
	if err != nil {
		srv.renderError(w, err)
		return
	}

	if n < p.Rounds {
		card, err := models.FindRandomCard(srv.Database, p.DeckID)
		if err != nil {
			srv.renderError(w, err)
			return
		}

		pr := &models.PracticeRound{
			PracticeID: p.MetaID,
			CardID:     card.MetaID,
			Expect:     card.Definitions[0], // FIXME cannot be first always
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
		pr, err := models.FindPracticeRound(srv.Database, p.MetaID, n)
		if err != nil {
			srv.renderError(w, err)
			return
		}

		card, err := models.FindCard(srv.Database, pr.CardID)
		if err != nil {
			srv.renderError(w, err)
			return
		}

		content.Card = card
		content.Round = pr
	}

	page = web.Page{
		Title:    "Practice",
		Partials: []string{"practice"},
		Content:  content,
	}

	srv.render(w, page)
}
