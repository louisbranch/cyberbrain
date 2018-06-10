package response

import (
	"context"
	"net/http"

	"gitlab.com/luizbranco/srs/web"
)

type Handler func(context.Context, http.ResponseWriter, *http.Request) Responder

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

func WrapError(err error, code int, msg string) Error {
	return Error{
		err:  err,
		code: code,
		msg:  msg,
	}
}

func NewError(code int, msg string) Error {
	return Error{
		code: code,
		msg:  msg,
	}
}

func (e Error) Error() string {
	return e.msg
}

func (e Error) FullError() string {
	if e.err != nil {
		return e.msg + ": " + e.err.Error()
	}

	return e.msg
}

func (e Error) Cause() error {
	return e.err
}

func (e Error) Code() int {
	return e.code
}

func (e Error) Respond(w http.ResponseWriter, r *http.Request) (*web.Page, error) {
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
