package response

import (
	"net/http"

	"github.com/luizbranco/srs/web"
)

type Handler func(http.ResponseWriter, *http.Request) Responder

type Responder interface {
	Respond(http.ResponseWriter, *http.Request) (*web.Page, error)
}

type Content struct {
	page web.Page
}

func (c Content) Respond(w http.ResponseWriter, r *http.Request) (*web.Page, error) {
	return &c.page, nil
}

func NewContent(page web.Page) Content {
	return Content{page: page}
}

type Error struct {
	err  error
	code int
	msg  string
}

func NewError(err error, code int, msg string) Error {
	return Error{
		err:  err,
		code: code,
		msg:  msg,
	}
}

func (e Error) Error() string {
	return e.msg
}

func (e Error) Respond(w http.ResponseWriter, r *http.Request) (*web.Page, error) {
	// TODO
	return nil, e
}

type Redirect struct {
	Path string
	Code int
}

func (rd Redirect) Respond(w http.ResponseWriter, r *http.Request) (*web.Page, error) {
	http.Redirect(w, r, rd.Path, rd.Code)
	return nil, nil
}
