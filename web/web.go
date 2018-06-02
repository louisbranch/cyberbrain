package web

import (
	"io"
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

type ID int

type Database interface {
	Create(Record) error
	Query(Query) ([]Record, error)
	QueryRaw(Query) ([]Record, error)
	Get(Query) (Record, error)
	Count(Query) (int, error)
	Random(Query, int) ([]Record, error)
}

type Record interface {
	ID() ID
	SetID(ID)
	Type() string
}

type Query interface {
	NewRecord() Record
	Where() map[string]interface{}
	Raw() string
}

type URLBuilder interface {
	Path(string, Record, ...Record) (string, error)
	ID(string) (ID, error)
}
