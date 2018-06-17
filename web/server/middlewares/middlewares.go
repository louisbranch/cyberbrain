package middlewares

import (
	"context"
	"net/http"

	"gitlab.com/luizbranco/cyberbrain/primitives"
	"gitlab.com/luizbranco/cyberbrain/web"
	"gitlab.com/luizbranco/cyberbrain/web/server/finder"
	"gitlab.com/luizbranco/cyberbrain/web/server/response"
)

type userKey struct{}
type deckKey struct{}

func Authenticate(h response.Handler) response.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) response.Responder {

		_, ok := CurrentUser(ctx)

		if !ok {
			return response.Redirect{Path: "/login", Code: http.StatusFound}
		}

		return h(ctx, w, r)
	}
}

func NewContext(user *primitives.User) context.Context {
	// TODO add request id

	ctx := context.Background()

	if user == nil {
		return ctx
	}

	u := *user
	u.PasswordHash = ""

	return context.WithValue(ctx, userKey{}, u)
}

func CurrentUser(ctx context.Context) (primitives.User, bool) {
	value := ctx.Value(userKey{})

	user, ok := value.(primitives.User)

	return user, ok
}

func Deck(h response.Handler, db primitives.Database, ub web.URLBuilder, hash string) response.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) response.Responder {

		user, _ := CurrentUser(ctx)

		deck, _, _, err := finder.Deck(db, ub, hash, finder.NoOption)
		if err != nil {
			return response.WrapError(err, http.StatusNotFound, "wrong deck id")
		}

		if user.ID() != deck.UserID {
			return response.WrapError(err, http.StatusForbidden, "invalid deck owner")
		}

		ctx = context.WithValue(ctx, deckKey{}, *deck)

		return h(ctx, w, r)
	}
}

func CurrentDeck(ctx context.Context) primitives.Deck {
	value := ctx.Value(deckKey{})

	return value.(primitives.Deck)
}
