package rounds

import (
	"net/http"

	"gitlab.com/luizbranco/srs"
	"gitlab.com/luizbranco/srs/db"
	"gitlab.com/luizbranco/srs/web"
	"gitlab.com/luizbranco/srs/web/html"
	"gitlab.com/luizbranco/srs/web/server/finder"
	"gitlab.com/luizbranco/srs/web/server/response"
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

		practice, err := finder.Practice(conn, ub, hash)
		if err != nil {
			return err.(response.Error)
		}

		// if practice is done, redirect back to its page
		if practice.Done {
			path, err := ub.Path("SHOW", practice)
			if err != nil {
				return response.WrapError(err, http.StatusInternalServerError, "failed to generate practice path")
			}

			return response.Redirect{Path: path, Code: http.StatusFound}
		}

		rounds, err := db.FindRounds(conn, practice.ID())
		if err != nil {
			return response.WrapError(err, http.StatusInternalServerError, "failed to find practice rounds")
		}

		// if a round is still in progress, redirect to its page
		for _, r := range rounds {
			if !r.Done {
				path, err := ub.Path("SHOW", r)
				if err != nil {
					return response.WrapError(err, http.StatusInternalServerError, "failed to generate round path")
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

		practice, err := finder.Practice(conn, ub, hash)
		if err != nil {
			return err.(response.Error)
		}

		// if practice is done, redirect back to its page
		if practice.Done {
			path, err := ub.Path("SHOW", practice)
			if err != nil {
				return response.WrapError(err, http.StatusInternalServerError, "failed to generate practice path")
			}

			return response.Redirect{Path: path, Code: http.StatusFound}
		}

		round, err := gen.NewRound(conn, *practice)
		if err != nil {
			return response.WrapError(err, http.StatusInternalServerError, "failed to generate new round")
		}

		path, err := ub.Path("SHOW", round)
		if err != nil {
			return response.WrapError(err, http.StatusInternalServerError, "failed to generate round path")
		}

		return response.Redirect{Path: path, Code: http.StatusFound}
	}
}

func Show(conn srs.Database, ub web.URLBuilder, hash string) response.Handler {
	return func(w http.ResponseWriter, r *http.Request) response.Responder {

		id, err := ub.ParseID(hash)
		if err != nil {
			return response.WrapError(err, http.StatusNotFound, "invalid round id")
		}

		round, err := db.FindRound(conn, id)
		if err != nil {
			return response.WrapError(err, http.StatusNotFound, "wrong round id")
		}

		practice, err := db.FindPractice(conn, round.PracticeID)
		if err != nil {
			return response.WrapError(err, http.StatusNotFound, "wrong practice id")
		}

		deck, err := db.FindDeck(conn, practice.DeckID)
		if err != nil {
			return response.WrapError(err, http.StatusInternalServerError, "invalid deck id")
		}

		round.Practice = practice
		practice.Deck = deck

		content, err := html.RenderRound(*round, ub)
		if err != nil {
			return response.WrapError(err, http.StatusInternalServerError, "failed to render round")
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
			return response.WrapError(err, http.StatusBadRequest, "invalid form")
		}

		id, err := ub.ParseID(hash)
		if err != nil {
			return response.WrapError(err, http.StatusNotFound, "invalid round id")
		}

		round, err := db.FindRound(conn, id)
		if err != nil {
			return response.WrapError(err, http.StatusNotFound, "wrong round id")
		}

		guess := r.Form.Get("guess")

		if guess == "" {
			return response.WrapError(err, http.StatusBadRequest, "invalid guess")
		}

		round.GuessAnswer(guess)

		err = conn.Update(round)
		if err != nil {
			return response.WrapError(err, http.StatusInternalServerError, "failed to update round")
		}

		practice, err := db.FindPractice(conn, round.PracticeID)
		if err != nil {
			return response.WrapError(err, http.StatusInternalServerError, "invalid practice id")
		}

		rounds, err := db.CountRounds(conn, practice.ID())
		if err != nil {
			return response.WrapError(err, http.StatusInternalServerError, "invalid rounds number")
		}

		if rounds >= practice.TotalRounds {
			practice.Done = true
			err = conn.Update(practice)
			if err != nil {
				return response.WrapError(err, http.StatusInternalServerError, "failed to update practice")
			}
		}

		path, err := ub.Path("SHOW", round)
		if err != nil {
			return response.WrapError(err, http.StatusInternalServerError, "failed to generate round path")
		}

		return response.Redirect{Path: path, Code: http.StatusFound}
	}
}
