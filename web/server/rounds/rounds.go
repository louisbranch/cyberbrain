package rounds

import (
	"net/http"

	"github.com/luizbranco/srs"
	"github.com/luizbranco/srs/db"
	"github.com/luizbranco/srs/web"
	"github.com/luizbranco/srs/web/html"
	"github.com/luizbranco/srs/web/server/response"
)

func Index() response.Handler {
	return func(w http.ResponseWriter, r *http.Request) response.Responder {
		return response.Redirect{Path: "/decks/", Code: http.StatusFound}
	}
}

func New(conn srs.Database, ub web.URLBuilder, gen srs.PracticeGenerator) response.Handler {
	return func(w http.ResponseWriter, r *http.Request) response.Responder {

		query := r.URL.Query()
		hash := query.Get("practice")

		practice, err := findPractice(conn, ub, hash)
		if err != nil {
			return err.(response.Error)
		}

		// if practice is done, redirect back to its page
		if practice.Done {
			path, err := ub.Path("SHOW", practice)
			if err != nil {
				return response.NewError(err, http.StatusInternalServerError, "failed to generate practice path")
			}

			return response.Redirect{Path: path, Code: http.StatusFound}
		}

		rounds, err := db.FindRounds(conn, practice.ID())
		if err != nil {
			return response.NewError(err, http.StatusInternalServerError, "failed to find practice rounds")
		}

		// if a round is still in progress, redirect to its page
		for _, r := range rounds {
			if !r.Done {
				path, err := ub.Path("SHOW", r)
				if err != nil {
					return response.NewError(err, http.StatusInternalServerError, "failed to generate round path")
				}

				return response.Redirect{Path: path, Code: http.StatusFound}
			}
		}

		// otherwise treat as a new round being created
		handler := Create(conn, ub, gen)
		return handler(w, r)
	}
}

func Create(conn srs.Database, ub web.URLBuilder, gen srs.PracticeGenerator) response.Handler {
	return func(w http.ResponseWriter, r *http.Request) response.Responder {

		query := r.URL.Query()
		hash := query.Get("practice")

		practice, err := findPractice(conn, ub, hash)
		if err != nil {
			return err.(response.Error)
		}

		round, err := gen.NewRound(conn, *practice)
		if err != nil {
			return response.NewError(err, http.StatusInternalServerError, "failed to generate new round")
		}

		path, err := ub.Path("SHOW", round)
		if err != nil {
			return response.NewError(err, http.StatusInternalServerError, "failed to generate round path")
		}

		return response.Redirect{Path: path, Code: http.StatusFound}
	}
}

func Show(conn srs.Database, ub web.URLBuilder, hash string) response.Handler {
	return func(w http.ResponseWriter, r *http.Request) response.Responder {

		id, err := ub.ParseID(hash)
		if err != nil {
			return response.NewError(err, http.StatusNotFound, "invalid round id")
		}

		round, err := db.FindRound(conn, id)
		if err != nil {
			return response.NewError(err, http.StatusNotFound, "wrong round id")
		}

		practice, err := db.FindPractice(conn, round.PracticeID)
		if err != nil {
			return response.NewError(err, http.StatusNotFound, "wrong practice id")
		}

		deck, err := db.FindDeck(conn, practice.DeckID)
		if err != nil {
			return response.NewError(err, http.StatusInternalServerError, "invalid deck id")
		}

		round.Practice = practice
		practice.Deck = deck

		content, err := html.RenderRound(*round, ub)
		if err != nil {
			return response.NewError(err, http.StatusInternalServerError, "failed to render round")
		}

		page := web.Page{
			Title:    "Practice Round",
			Partials: []string{"round"},
			Content:  content,
		}

		return response.NewContent(page)
	}
}

func Update(conn srs.Database, ub web.URLBuilder, hash string) response.Handler {
	return func(w http.ResponseWriter, r *http.Request) response.Responder {
		if err := r.ParseForm(); err != nil {
			return response.NewError(err, http.StatusBadRequest, "invalid form")
		}

		id, err := ub.ParseID(hash)
		if err != nil {
			return response.NewError(err, http.StatusNotFound, "invalid round id")
		}

		round, err := db.FindRound(conn, id)
		if err != nil {
			return response.NewError(err, http.StatusNotFound, "wrong round id")
		}

		guess := r.Form.Get("guess")

		if guess == "" {
			return response.NewError(err, http.StatusBadRequest, "invalid guess")
		}

		round.GuessAnswer(guess)

		err = conn.Update(round)
		if err != nil {
			return response.NewError(err, http.StatusInternalServerError, "failed to update round")
		}

		path, err := ub.Path("SHOW", round)
		if err != nil {
			return response.NewError(err, http.StatusInternalServerError, "failed to generate round path")
		}

		return response.Redirect{Path: path, Code: http.StatusFound}
	}
}

func findPractice(conn srs.Database, ub web.URLBuilder, hash string) (*srs.Practice, error) {
	id, err := ub.ParseID(hash)
	if err != nil {
		return nil, response.NewError(err, http.StatusBadRequest, "invalid practice id")
	}

	practice, err := db.FindPractice(conn, id)
	if err != nil {
		return nil, response.NewError(err, http.StatusBadRequest, "wrong practice id")
	}

	return practice, nil
}

func findDeck(conn srs.Database, ub web.URLBuilder, hash string) (*srs.Deck, error) {

	id, err := ub.ParseID(hash)
	if err != nil {
		return nil, response.NewError(err, http.StatusBadRequest, "invalid deck id")
	}

	deck, err := db.FindDeck(conn, id)
	if err != nil {
		return nil, response.NewError(err, http.StatusBadRequest, "wrong deck id")
	}

	tags, err := db.FindTags(conn, id)
	if err != nil {
		return nil, response.NewError(err, http.StatusInternalServerError, "failed to find deck tags")
	}

	deck.Tags = tags

	return deck, nil
}
