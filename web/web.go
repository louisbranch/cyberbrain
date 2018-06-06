package web

import (
	"io"

	"gitlab.com/luizbranco/srs/primitives"
)

type Page struct {
	Title      string
	ActiveMenu string
	Layout     string
	Partials   []string
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
