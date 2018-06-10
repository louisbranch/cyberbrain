package middlewares

import (
	"context"
	"net/http"

	"gitlab.com/luizbranco/srs/primitives"
	"gitlab.com/luizbranco/srs/web/server/response"
)

var userKey struct{}

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

	return context.WithValue(ctx, userKey, u)
}

func CurrentUser(ctx context.Context) (primitives.User, bool) {
	valeu := ctx.Value(userKey)

	user, ok := valeu.(primitives.User)

	return user, ok
}
