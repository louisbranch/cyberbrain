package web

import (
	"io"

	"github.com/luizbranco/srs"
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
	ParseID(string) (srs.ID, error)
	EncodeID(srs.ID) (string, error)

	Path(string, srs.Identifiable, ...srs.Identifiable) (string, error)
}
