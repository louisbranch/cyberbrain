package web

import (
	"database/sql/driver"
	"io"
	"strconv"

	"github.com/pkg/errors"
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

type ID string

func (id *ID) Value() (driver.Value, error) {
	return strconv.ParseInt(string(*id), 10, 64)
}

func (id *ID) Scan(v interface{}) error {
	switch t := v.(type) {
	case int64:
		sid := ID(strconv.FormatInt(v.(int64), 10))
		*id = sid
		return nil
	default:
		return errors.Errorf("failed to scan value %q into ID", t)
	}
}

type Database interface {
	Create(Record) error
	Query(Query) ([]Record, error)
	QueryRaw(Query) ([]Record, error)
	Get(Query) (Record, error)
	Count(Query) (int, error)
	Random(Query, int) ([]Record, error)
}

type Record interface {
	SetID(ID)
	Type() string
}

type Query interface {
	NewRecord() Record
	Where() map[string]interface{}
	Raw() string
}
