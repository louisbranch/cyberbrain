package server

import (
	"net/http"

	"github.com/luizbranco/srs/web"
	"github.com/luizbranco/srs/web/models"
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
			srv.renderError(w, err)
			return
		}

		p, err := models.NewPracticeFromForm(deck.ID, r.Form)
		if err != nil {
			srv.renderError(w, err)
			return
		}

		err = srv.Database.Create(p)
		if err != nil {
			srv.renderError(w, err)
			return
		}

		http.Redirect(w, r, "/practices/"+p.Slug, http.StatusFound)
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

	deck, err := models.FindDeckBySlug(srv.Database, slug)
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
	p, err := models.FindPracticeBySlug(srv.Database, slug)
	if err != nil {
		srv.renderError(w, err)
		return
	}

	deck, err := models.FindDeckByID(srv.Database, p.DeckID)
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
		Title:    "Practice #" + slug,
		Partials: []string{"practice"},
		Content:  content,
	}

	if p.Finished() {
		srv.render(w, page)
		return
	}

	n, err := models.CountPracticeRounds(srv.Database, p.ID)
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
		pr, err := models.FindPracticeRound(srv.Database, p.ID, n)
		if err != nil {
			srv.renderError(w, err)
			return
		}

		card, err := models.FindCardByID(srv.Database, pr.CardID)
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
