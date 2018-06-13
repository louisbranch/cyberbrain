package psql

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/lib/pq"
	"github.com/pkg/errors"
	"gitlab.com/luizbranco/cyberbrain/primitives"
)

type Query struct {
	table   string
	columns []string
	fields  []interface{}
	addrs   []interface{}
}

func QueryFromRecord(r primitives.Record, ignored ...string) (*Query, error) {
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

		field := rv.Field(i)
		addr := field.Addr().Interface()

		q.fields = append(q.fields, addr)

		if field.Kind() == reflect.Slice {
			switch e := field.Type().Elem(); e.Kind() {
			case reflect.String:
				slice := field.Interface().([]string)
				sa := pq.StringArray(slice)
				addr = &sa
			case reflect.Int:
				ids := field.Interface().([]primitives.ID)
				slice := make([]int64, len(ids))

				for i, id := range ids {
					slice[i] = int64(id)
				}

				sa := pq.Int64Array(slice)
				addr = &sa
			case reflect.Uint8: // []byte
			default:
				return nil, errors.Errorf("slice type %v not supported", e)
			}
		}

		q.addrs = append(q.addrs, addr)
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

type Scannable interface {
	Scan(...interface{}) error
}

func (q *Query) Scan(row Scannable) error {
	err := row.Scan(q.addrs...)
	if err != nil {
		return errors.Wrap(err, "failed to scan records")
	}

	for i, addr := range q.addrs {
		switch addr.(type) {
		case *pq.StringArray:
			slice := addr.(*pq.StringArray)
			ss := []string(*slice)
			f := q.fields[i].(*[]string)
			*f = ss
		case *pq.Int64Array:
			slice := addr.(*pq.Int64Array)
			f := q.fields[i].(*[]primitives.ID)

			sids := make([]primitives.ID, len(*slice))
			for i, id := range *slice {
				sids[i] = primitives.ID(id)
			}

			*f = sids
		}
	}

	return nil
}

func where(cond primitives.Query) string {
	where := cond.Where()

	if len(where) == 0 {
		return ""
	}

	var clause []string
	for k, v := range where {
		switch t := v.(type) {
		case string:
			clause = append(clause, fmt.Sprintf("%s = '%s'", k, v))
		case int, primitives.ID:
			clause = append(clause, fmt.Sprintf("%s = %d", k, v))
		default:
			err := fmt.Sprintf("invalid type %q for where clause", t)
			panic(err)
		}
	}

	q := "WHERE " + strings.Join(clause, " AND ")

	sortBy := cond.SortBy()

	if len(sortBy) == 0 {
		return q
	}

	var sort []string

	for k, v := range sortBy {
		sort = append(sort, fmt.Sprintf("%s %s", k, v))
	}

	return q + "ORDER BY " + strings.Join(sort, " AND ")
}
