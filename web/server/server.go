package server

import (
	"fmt"
	"net/http"

	"github.com/luizbranco/srs/web"
)

type Server struct {
	Template web.Template
	Database web.Database
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

	mux.HandleFunc("/practices/new", srv.newPractice)
	mux.HandleFunc("/practices/", srv.practice)

	mux.HandleFunc("/", srv.index)

	return mux
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
