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
	Query(string, Collection) error
	Get(uint, Record) error
}

type Record interface {
	SetID(uint)
	Type() string
}

type Collection interface {
	NewRecord() Record
	Append(Record)
}
