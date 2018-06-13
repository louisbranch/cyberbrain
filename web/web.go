package web

import (
	"io"
	"net/http"

	"gitlab.com/luizbranco/cyberbrain/primitives"
)

type Page struct {
	Title      string
	ActiveMenu string
	Layout     string
	Partials   []string
	User       *primitives.User
	Content    interface{}
}

type Template interface {
	Render(w io.Writer, page Page) error
}

type URLBuilder interface {
	ParseID(string) (primitives.ID, error)
	EncodeID(primitives.ID) (string, error)

	Path(string, primitives.Identifiable, ...primitives.Identifiable) (string, error)
}

type SessionManager interface {
	LogIn(primitives.User, http.ResponseWriter) error
	LogOut(http.ResponseWriter)
	User(*http.Request) (*primitives.User, error)
}
