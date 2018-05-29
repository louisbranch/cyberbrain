package sqlite

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
		table: r.Type(),
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

func (q *Query) Table() string {
	return q.table
}

func (q *Query) Placeholders() string {
	v := make([]string, len(q.columns))
	for i := range v {
		v[i] = "?"
	}

	return strings.Join(v, ", ")
}

func (q *Query) Columns() string {
	return strings.Join(q.columns, ", ")
}

func contains(list []string, s string) bool {
	for i := range list {
		if list[i] == s {
			return true
		}
	}
	return false
}

func where(cond web.Condition) string {
	if len(cond.Where) == 0 {
		return ""
	}

	var clause []string
	for k, v := range cond.Where {
		switch t := v.(type) {
		case string:
			clause = append(clause, fmt.Sprintf("%s = %q", k, v))
		case uint, uint64, int, int64, int32:
			clause = append(clause, fmt.Sprintf("%s = %d", k, v))
		default:
			err := fmt.Sprintf("invalid type %v for where clause", t)
			panic(err)
		}
	}

	return "WHERE " + strings.Join(clause, " AND ")
}
