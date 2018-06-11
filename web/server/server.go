package server

import (
	"fmt"
	"log"
	"net/http"

	"gitlab.com/luizbranco/srs/primitives"
	"gitlab.com/luizbranco/srs/web"
	"gitlab.com/luizbranco/srs/web/server/decks"
	"gitlab.com/luizbranco/srs/web/server/home"
	"gitlab.com/luizbranco/srs/web/server/middlewares"
	"gitlab.com/luizbranco/srs/web/server/response"
	"gitlab.com/luizbranco/srs/web/server/sessions"
	"gitlab.com/luizbranco/srs/web/server/users"
	"gitlab.com/luizbranco/srs/worker/jobs"
)

type Server struct {
	Template          web.Template
	URLBuilder        web.URLBuilder
	Database          primitives.Database
	PracticeGenerator primitives.PracticeGenerator
	Authenticator     primitives.Authenticator
	SessionManager    web.SessionManager
	ImageUploader     jobs.ImageUploader
}

func (srv *Server) NewServeMux() *http.ServeMux {
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("web/assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))

	renderer := &middlewares.Renderer{
		SessionManager: srv.SessionManager,
		Template:       srv.Template,
	}

	signupMux := users.NewSignupMux(renderer, srv.Database, srv.URLBuilder,
		srv.Authenticator)

	loginMux := sessions.NewLoginMux(renderer, srv.Database, srv.URLBuilder,
		srv.Authenticator)

	logoutMux := sessions.NewLogoutMux(renderer)

	decksMux := decks.NewServeMux(renderer, srv.Database, srv.URLBuilder,
		srv.PracticeGenerator, srv.ImageUploader)

	mux.Handle("/signup/", http.StripPrefix("/signup", signupMux))
	mux.Handle("/login/", http.StripPrefix("/login", loginMux))
	mux.Handle("/logout/", http.StripPrefix("/logout", logoutMux))
	mux.Handle("/decks/", http.StripPrefix("/decks", decksMux))
	mux.Handle("/", home.NewServeMux(renderer))

	return mux
}

func (srv *Server) handle(handler response.Handler, w http.ResponseWriter, r *http.Request) {
	if handler == nil {
		err := response.NewError(http.StatusNotFound, r.URL.Path+" not found")
		srv.renderError(w, err, nil)
		return
	}

	user, err := srv.SessionManager.User(r)
	if err != nil {
		srv.renderError(w, err, nil)
		return
	}

	ctx := middlewares.NewContext(user)

	res := handler(ctx, w, r)
	page, err := res.Respond(w, r)

	if page != nil {
		page.User = user
		srv.render(w, *page)
		return
	}

	if err != nil {
		srv.renderError(w, err, user)
		return
	}
}

func (srv *Server) render(w http.ResponseWriter, page web.Page) {
	if page.Layout == "" {
		page.Layout = "layout"
	}

	err := srv.Template.Render(w, page)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err)
	}
}

func (srv *Server) renderError(w http.ResponseWriter, err error, user *primitives.User) {
	code := http.StatusInternalServerError

	res, ok := err.(response.Error)

	if ok {
		code = res.Code()
	}

	var page web.Page

	switch code {
	case http.StatusNotFound:
		page = web.Page{
			Title:    "Not Found",
			Partials: []string{"404"},
			User:     user,
		}
	case http.StatusBadRequest:
		page = web.Page{
			Title:    "Bad Request",
			Content:  err,
			Partials: []string{"400"},
			User:     user,
		}
	default:
		code = http.StatusInternalServerError
		page = web.Page{
			Title:    "Internal Server Error",
			Content:  err,
			Partials: []string{"500"},
			User:     user,
		}
	}

	log.Println(err.Error())

	w.WriteHeader(code)

	srv.render(w, page)
}
