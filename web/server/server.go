package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"gitlab.com/luizbranco/srs/primitives"
	"gitlab.com/luizbranco/srs/web"
	"gitlab.com/luizbranco/srs/web/server/cards"
	"gitlab.com/luizbranco/srs/web/server/decks"
	"gitlab.com/luizbranco/srs/web/server/middlewares"
	"gitlab.com/luizbranco/srs/web/server/practices"
	"gitlab.com/luizbranco/srs/web/server/response"
	"gitlab.com/luizbranco/srs/web/server/rounds"
	"gitlab.com/luizbranco/srs/web/server/sessions"
	"gitlab.com/luizbranco/srs/web/server/tags"
	"gitlab.com/luizbranco/srs/web/server/users"
)

type Server struct {
	Template          web.Template
	URLBuilder        web.URLBuilder
	Database          primitives.Database
	PracticeGenerator primitives.PracticeGenerator
	Authenticator     primitives.Authenticator
	SessionManager    web.SessionManager
}

func (srv *Server) NewServeMux() *http.ServeMux {
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("web/assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))

	mux.HandleFunc("/signup/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path[len("/signup/"):]
		method := r.Method

		var handler response.Handler

		switch {
		case method == "GET" && path == "":
			handler = users.New(srv.Database, srv.URLBuilder)
		case method == "POST" && path == "":
			handler = users.Create(srv.Database, srv.URLBuilder, srv.Authenticator, srv.SessionManager)
		}

		srv.handle(handler, w, r)
	})

	mux.HandleFunc("/login/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path[len("/login/"):]
		method := r.Method

		var handler response.Handler

		switch {
		case method == "GET" && path == "":
			handler = sessions.New(srv.Database, srv.URLBuilder)
		case method == "POST" && path == "":
			handler = sessions.Create(srv.Database, srv.URLBuilder, srv.Authenticator, srv.SessionManager)
		}

		srv.handle(handler, w, r)
	})

	mux.HandleFunc("/logout/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path[len("/logout/"):]
		method := r.Method

		var handler response.Handler

		switch {
		case method == "GET" && path == "":
			handler = sessions.Destroy(srv.SessionManager)
		}

		srv.handle(handler, w, r)
	})

	mux.HandleFunc("/decks/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path[len("/decks/"):]
		method := r.Method

		var handler response.Handler

		switch {
		case method == "GET" && path == "":
			handler = decks.Index(srv.Database, srv.URLBuilder)
		case method == "GET" && path == "new":
			handler = decks.New(srv.Database, srv.URLBuilder)
		case method == "GET":
			handler = decks.Show(srv.Database, srv.URLBuilder, path)
		case method == "POST" && path == "":
			handler = decks.Create(srv.Database, srv.URLBuilder)
		}

		handler = middlewares.Authenticate(handler)

		srv.handle(handler, w, r)
	})

	mux.HandleFunc("/cards/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path[len("/cards/"):]
		method := r.Method

		var handler response.Handler

		switch {
		case method == "GET" && path == "":
			handler = cards.Index()
		case method == "GET" && path == "new":
			handler = cards.New(srv.Database, srv.URLBuilder)
		case method == "GET":
			handler = cards.Show(srv.Database, srv.URLBuilder, path)
		case method == "POST" && path == "":
			handler = cards.Create(srv.Database, srv.URLBuilder)
		case method == "POST":
			handler = cards.Update(srv.Database, srv.URLBuilder, path)
		}

		handler = middlewares.Authenticate(handler)

		srv.handle(handler, w, r)
	})

	mux.HandleFunc("/tags/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path[len("/tags/"):]
		method := r.Method

		var handler response.Handler

		switch {
		case method == "GET" && path == "":
			handler = tags.Index()
		case method == "GET" && path == "new":
			handler = tags.New(srv.Database, srv.URLBuilder)
		case method == "GET":
			handler = tags.Show(srv.Database, srv.URLBuilder, path)
		case method == "POST" && path == "":
			handler = tags.Create(srv.Database, srv.URLBuilder)
		case method == "POST":
			handler = tags.Update(srv.Database, srv.URLBuilder, path)
		}

		handler = middlewares.Authenticate(handler)

		srv.handle(handler, w, r)
	})

	mux.HandleFunc("/practices/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path[len("/practices/"):]
		method := r.Method

		var handler response.Handler

		switch {
		case method == "GET" && path == "":
			handler = practices.Index()
		case method == "GET" && path == "new":
			handler = practices.New(srv.Database, srv.URLBuilder)
		case method == "GET":
			handler = practices.Show(srv.Database, srv.URLBuilder, path)
		case method == "POST" && path == "":
			handler = practices.Create(srv.Database, srv.URLBuilder)
		}

		handler = middlewares.Authenticate(handler)

		srv.handle(handler, w, r)
	})

	mux.HandleFunc("/rounds/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path[len("/rounds/"):]
		method := r.Method

		var handler response.Handler

		switch {
		case method == "GET" && path == "":
			handler = rounds.Index()
		case method == "GET" && path == "new":
			handler = rounds.New(srv.Database, srv.URLBuilder, srv.PracticeGenerator)
		case method == "GET":
			handler = rounds.Show(srv.Database, srv.URLBuilder, path)
		case method == "POST" && path == "":
			handler = rounds.Create(srv.Database, srv.URLBuilder, srv.PracticeGenerator)
		case method == "POST":
			handler = rounds.Update(srv.Database, srv.URLBuilder, path)
		}

		handler = middlewares.Authenticate(handler)

		srv.handle(handler, w, r)
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		srv.handle(index(), w, r)
	})

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

func index() response.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) response.Responder {
		if r.Method != "GET" || r.URL.Path != "/" {
			return response.NewError(http.StatusNotFound, r.URL.Path+" not found")
		}

		return response.Redirect{Path: "/decks/", Code: http.StatusFound}
	}
}
