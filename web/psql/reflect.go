package psql

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/luizbranco/srs/web"
	"github.com/pkg/errors"
)

type Query struct {
	table   string
	columns []string
	fields  []interface{}
}

func QueryFromRecord(r web.Record, ignored ...string) (*Query, error) {
	rv := reflect.ValueOf(r)
	if rv.Kind() != reflect.Ptr {
		return nil, errors.Errorf("cannot get database fields for record %v", r)
	}

	q := &Query{
		table: r.Type() + "s",
	}

	rv = rv.Elem()
	rt := reflect.TypeOf(rv.Interface())

	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		tag := f.Tag.Get("db")

		if tag == "" || contains(ignored, tag) {
			continue
		}

		field := rv.Field(i).Addr().Interface()
		q.fields = append(q.fields, field)
		q.columns = append(q.columns, tag)
	}

	return q, nil
}

func contains(l []string, s string) bool {
	for _, i := range l {
		if i == s {
			return true
		}
	}
	return false
}

func (q *Query) Table() string {
	return q.table
}

func (q *Query) Placeholders() string {
	v := make([]string, len(q.columns))
	for i := range v {
		v[i] = fmt.Sprintf("$%d", i+1)
	}

	return strings.Join(v, ", ")
}

func (q *Query) Columns() string {
	return strings.Join(q.columns, ", ")
}

func where(cond web.Query) string {
	where := cond.Where()

	if len(where) == 0 {
		return ""
	}

	var clause []string
	for k, v := range where {
		switch t := v.(type) {
		case string, web.ID:
			clause = append(clause, fmt.Sprintf("%s = '%s'", k, v))
		case int:
			clause = append(clause, fmt.Sprintf("%s = %d", k, v))
		default:
			err := fmt.Sprintf("invalid type %q for where clause", t)
			panic(err)
		}
	}

	return "WHERE " + strings.Join(clause, " AND ")
}
