package reviews

import (
	"context"
	"net/http"

	"gitlab.com/luizbranco/cyberbrain/db"
	"gitlab.com/luizbranco/cyberbrain/primitives"
	"gitlab.com/luizbranco/cyberbrain/web"
	"gitlab.com/luizbranco/cyberbrain/web/html"
	"gitlab.com/luizbranco/cyberbrain/web/server/finder"
	"gitlab.com/luizbranco/cyberbrain/web/server/middlewares"
	"gitlab.com/luizbranco/cyberbrain/web/server/response"
)

func Index() response.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) response.Responder {
		return response.Redirect{Path: "/decks/", Code: http.StatusFound}
	}
}

func New(conn primitives.Database, ub web.URLBuilder) response.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) response.Responder {

		deck := middlewares.CurrentDeck(ctx)

		schedule, err := db.FindNextCardScheduled(conn, deck.ID())
		if err != nil {
			handler := Summary(conn, ub)
			return handler(ctx, w, r)
		}

		card, _, err := finder.Card(conn, ub, schedule.CardID, finder.NoOption)
		if err != nil {
			return err.(response.Error)
		}

		cardC, err := html.RenderCard(ub, deck, nil, *card, nil, true)
		if err != nil {
			return response.WrapError(err, http.StatusInternalServerError, "failed to render card")
		}

		path, err := ub.Path("CREATE", &primitives.CardReview{}, deck)
		if err != nil {
			return response.WrapError(err, http.StatusInternalServerError, "failed to build deck create review path")
		}

		content := struct {
			Card       *html.Card
			ReviewPath string
		}{
			Card:       cardC,
			ReviewPath: path,
		}

		page := web.Page{
			Title:    "Card Review",
			Partials: []string{"review"},
			Content:  content,
		}

		return response.NewContent(page)
	}
}

func Summary(conn primitives.Database, ub web.URLBuilder) response.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) response.Responder {

		deck := middlewares.CurrentDeck(ctx)

		deckC, err := html.RenderDeck(ub, deck, nil, nil)
		if err != nil {
			return response.WrapError(err, http.StatusInternalServerError, "failed to render deck")
		}

		content := struct {
			Deck *html.Deck
		}{
			Deck: deckC,
		}

		page := web.Page{
			Title:    "New Review",
			Partials: []string{"empty_review"},
			Content:  content,
		}

		return response.NewContent(page)
	}
}

func Create(conn primitives.Database, ub web.URLBuilder) response.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) response.Responder {

		if err := r.ParseForm(); err != nil {
			return response.WrapError(err, http.StatusBadRequest, "invalid form")
		}

		deck := middlewares.CurrentDeck(ctx)

		cardID := r.Form.Get("card_id")

		card, _, err := finder.Card(conn, ub, cardID, finder.NoOption)
		if err != nil {
			return err.(response.Error)
		}

		review := &primitives.CardReview{
			DeckID: deck.ID(),
			CardID: card.ID(),
			Answer: r.Form.Get("answer"),
		}

		action := r.Form.Get("action")
		if action == "Skip" {
			review.Skipped = true
		}

		for _, d := range card.Definitions {
			if d == review.Answer {
				review.Correct = true
				break
			}
		}

		err = conn.Create(review)
		if err != nil {
			return response.WrapError(err, http.StatusInternalServerError, "failed to create card review")
		}

		schedule, err := db.FindCardSchedule(conn, card.ID())
		if err != nil {
			schedule = primitives.NewCardSchedule(deck.ID(), card.ID())

			err := conn.Create(schedule)
			if err != nil {
				return response.WrapError(err, http.StatusInternalServerError, "failed to create card schedule")
			}
		}

		schedule.Reschedule(review.Correct)

		err = conn.Update(schedule)
		if err != nil {
			return response.WrapError(err, http.StatusInternalServerError, "failed to update card schedule")
		}

		path, err := ub.Path("SHOW", review, deck)
		if err != nil {
			return response.WrapError(err, http.StatusInternalServerError, "failed to build deck show review path")
		}

		return response.Redirect{Path: path, Code: http.StatusFound}
	}
}

func Show(conn primitives.Database, ub web.URLBuilder, hash string) response.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) response.Responder {

		deck := middlewares.CurrentDeck(ctx)

		review, err := finder.CardReview(conn, ub, hash)
		if err != nil {
			return err.(response.Error)
		}

		card, _, err := finder.Card(conn, ub, review.CardID, finder.NoOption)
		if err != nil {
			return err.(response.Error)
		}

		cardC, err := html.RenderCard(ub, deck, nil, *card, nil, true)
		if err != nil {
			return response.WrapError(err, http.StatusInternalServerError, "failed to render card")
		}

		content := struct {
			Card   *html.Card
			Review *primitives.CardReview
		}{
			Card:   cardC,
			Review: review,
		}

		page := web.Page{
			Title:    "Card Review Result",
			Partials: []string{"review_result"},
			Content:  content,
		}

		return response.NewContent(page)
	}
}
