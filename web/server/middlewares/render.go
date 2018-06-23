package middlewares

import (
	"log"
	"net/http"

	"gitlab.com/luizbranco/cyberbrain/primitives"
	"gitlab.com/luizbranco/cyberbrain/web"
	"gitlab.com/luizbranco/cyberbrain/web/server/response"
)

type Renderer struct {
	SessionManager web.SessionManager
	Template       web.Template
}

func (rr *Renderer) Render(handler response.Handler, w http.ResponseWriter, r *http.Request) {
	log.Printf("[%s] %s\n", r.Method, r.RequestURI)

	if handler == nil {
		err := response.NewError(http.StatusNotFound, r.URL.Path+" not found")
		rr.renderError(w, r, err, nil)
		return
	}

	user, err := rr.SessionManager.User(r)
	if err != nil {
		rr.SessionManager.LogOut(w)
		log.Println(err)
	}

	ctx := NewContext(user)

	res := handler(ctx, w, r)
	page, err := res.Respond(w, r)

	if page != nil {
		page.User = user
		err = rr.render(w, *page)
	}

	if err != nil {
		rr.renderError(w, r, err, user)
		return
	}
}

func (rr *Renderer) render(w http.ResponseWriter, page web.Page) error {
	if page.Layout == "" {
		page.Layout = "layout"
	}

	err := rr.Template.Render(w, page)
	if err != nil {
		return err
	}

	return nil
}

func (rr *Renderer) renderError(w http.ResponseWriter, r *http.Request, err error, user *primitives.User) {
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

	switch err := err.(type) {
	case response.Error:
		log.Printf("ERROR - [%s] %s (%s)\n", r.Method, r.RequestURI, err.FullError())
	default:
		log.Printf("ERROR - [%s] %s (%s)\n", r.Method, r.RequestURI, err.Error())
	}

	w.WriteHeader(code)

	rr.render(w, page)
}
