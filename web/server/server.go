package server

import (
	"net/http"

	"gitlab.com/luizbranco/cyberbrain/primitives"
	"gitlab.com/luizbranco/cyberbrain/web"
	"gitlab.com/luizbranco/cyberbrain/web/server/blitline"
	"gitlab.com/luizbranco/cyberbrain/web/server/decks"
	"gitlab.com/luizbranco/cyberbrain/web/server/home"
	"gitlab.com/luizbranco/cyberbrain/web/server/middlewares"
	"gitlab.com/luizbranco/cyberbrain/web/server/sessions"
	"gitlab.com/luizbranco/cyberbrain/web/server/users"
	"gitlab.com/luizbranco/cyberbrain/worker"
)

type Server struct {
	Template       web.Template
	URLBuilder     web.URLBuilder
	Database       primitives.Database
	Authenticator  primitives.Authenticator
	SessionManager web.SessionManager
	ImageResizer   worker.ImageResizer
}

func (srv *Server) NewServeMux() *http.ServeMux {
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("web/assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))

	renderer := &middlewares.Renderer{
		SessionManager: srv.SessionManager,
		Template:       srv.Template,
	}

	signupMux := users.NewServeMux(renderer, srv.Database, srv.URLBuilder,
		srv.Authenticator)

	loginMux := sessions.NewLoginMux(renderer, srv.Database, srv.URLBuilder,
		srv.Authenticator)

	logoutMux := sessions.NewLogoutMux(renderer)

	decksMux := decks.NewServeMux(renderer, srv.Database, srv.URLBuilder,
		srv.ImageResizer)

	blitlineMux := blitline.NewServeMux(renderer, srv.Database, srv.URLBuilder)

	mux.Handle("/signup/", http.StripPrefix("/signup", signupMux))
	mux.Handle("/login/", http.StripPrefix("/login", loginMux))
	mux.Handle("/logout/", http.StripPrefix("/logout", logoutMux))
	mux.Handle("/decks/", http.StripPrefix("/decks", decksMux))
	mux.Handle("/blitline/", http.StripPrefix("/blitline", blitlineMux))

	mux.HandleFunc("/_healthz/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mux.Handle("/", home.NewServeMux(renderer))

	return mux
}
