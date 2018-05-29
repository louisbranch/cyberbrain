package web

import "io"

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

type Database interface {
	Create(Record) error
	Query(Condition) ([]Record, error)
	QueryRaw(Condition) ([]Record, error)
	Get(Condition) (Record, error)
	Count(Condition) (int, error)
	Random(Condition, int) ([]Record, error)
}

type Record interface {
	SetID(uint)
	Type() string
	GenerateSlug() error
}

type Condition struct {
	Record Record
	Where  map[string]interface{}
	Raw    string
}
