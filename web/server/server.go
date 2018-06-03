package server

import (
	"fmt"
	"net/http"

	"github.com/luizbranco/srs"
	"github.com/luizbranco/srs/web"
	"github.com/luizbranco/srs/web/server/practices"
	"github.com/luizbranco/srs/web/server/response"
	"github.com/luizbranco/srs/web/server/rounds"
)

type Server struct {
	Template          web.Template
	URLBuilder        web.URLBuilder
	Database          srs.Database
	PracticeGenerator srs.PracticeGenerator
}

func (srv *Server) NewServeMux() *http.ServeMux {
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("web/assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))

	mux.HandleFunc("/decks/new/", srv.newDeck)
	mux.HandleFunc("/decks/", srv.decks)

	mux.HandleFunc("/cards/new", srv.newCard)
	mux.HandleFunc("/cards/", srv.cards)

	mux.HandleFunc("/tags/new", srv.newTag)
	mux.HandleFunc("/tags/", srv.tags)

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
			handler = rounds.Create(srv.Database, srv.URLBuilder)
		}

		srv.handle(handler, w, r)
	})

	mux.HandleFunc("/", srv.index)

	return mux
}

func (srv *Server) handle(handler response.Handler, w http.ResponseWriter, r *http.Request) {
	if handler == nil {
		srv.renderNotFound(w)
		return
	}

	res := handler(w, r)
	page, err := res.Respond(w, r)

	if page != nil {
		srv.render(w, *page)
		return
	}

	if err != nil {
		srv.renderError(w, err)
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

func (srv *Server) renderError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	page := web.Page{
		Title:    "500",
		Content:  err,
		Partials: []string{"500"},
	}
	srv.render(w, page)
}

func (srv *Server) renderNotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	page := web.Page{
		Title:    "Not Found",
		Partials: []string{"404"},
	}
	srv.render(w, page)
}

func (srv *Server) index(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/" {
		srv.renderNotFound(w)
		return
	}

	http.Redirect(w, r, "/decks/", http.StatusMovedPermanently)
}
