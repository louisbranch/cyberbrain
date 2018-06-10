package middlewares

import (
	"fmt"
	"log"
	"net/http"

	"gitlab.com/luizbranco/srs/primitives"
	"gitlab.com/luizbranco/srs/web"
	"gitlab.com/luizbranco/srs/web/server/response"
)

type Renderer struct {
	SessionManager web.SessionManager
	Template       web.Template
}

func (rr *Renderer) Render(handler response.Handler, w http.ResponseWriter, r *http.Request) {
	if handler == nil {
		err := response.NewError(http.StatusNotFound, r.URL.Path+" not found")
		rr.renderError(w, err, nil)
		return
	}

	user, err := rr.SessionManager.User(r)
	if err != nil {
		rr.renderError(w, err, nil)
		return
	}

	ctx := NewContext(user)

	res := handler(ctx, w, r)
	page, err := res.Respond(w, r)

	if page != nil {
		page.User = user
		rr.render(w, *page)
		return
	}

	if err != nil {
		rr.renderError(w, err, user)
		return
	}
}

func (rr *Renderer) render(w http.ResponseWriter, page web.Page) {
	if page.Layout == "" {
		page.Layout = "layout"
	}

	err := rr.Template.Render(w, page)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err)
	}
}

func (rr *Renderer) renderError(w http.ResponseWriter, err error, user *primitives.User) {
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

	rr.render(w, page)
}
