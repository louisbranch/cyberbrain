package sqlite

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/luizbranco/srs/web"
	"github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

type Database struct {
	*sql.DB
}

func init() {
	sql.Register("sqlite3_with_fk",
		&sqlite3.SQLiteDriver{
			ConnectHook: func(conn *sqlite3.SQLiteConn) error {
				_, err := conn.Exec("PRAGMA foreign_keys = ON", nil)
				return err
			},
		})
}

func New(path string) (*Database, error) {
	db, err := sql.Open("sqlite3_with_fk", path)
	if err != nil {
		return nil, err
	}

	queries := []string{
		`
		CREATE TABLE IF NOT EXISTS decks(
			id INTEGER PRIMARY KEY,
			slug TEXT NOT NULL UNIQUE CHECK(slug <> ''),
			name TEXT NOT NULL CHECK(name <> ''),
			description TEXT,
			image_url TEXT,
			field_1 TEXT NOT NULL CHECK(field_1 <> ''),
			field_2 TEXT NOT NULL CHECK(field_2 <> ''),
			field_3 TEXT NOT NULL CHECK(field_3 <> '')
		);
		`,
		`
		CREATE TABLE IF NOT EXISTS cards(
			id INTEGER PRIMARY KEY,
			slug TEXT NOT NULL UNIQUE CHECK(slug <> ''),
			deck_id INTEGER NOT NULL,
			image_url TEXT,
			audio_url TEXT,
			field_1 TEXT NOT NULL CHECK(field_1 <> ''),
			field_2 TEXT NOT NULL CHECK(field_2 <> ''),
			field_3 TEXT NOT NULL CHECK(field_3 <> ''),
			FOREIGN KEY(deck_id) REFERENCES decks(id) ON DELETE CASCADE
		);
		`,
		`
		CREATE TABLE IF NOT EXISTS tags(
			id INTEGER PRIMARY KEY,
			slug TEXT NOT NULL UNIQUE CHECK(slug <> ''),
			name TEXT NOT NULL UNIQUE CHECK(name <> ''),
			deck_id INTEGER NOT NULL,
			FOREIGN KEY(deck_id) REFERENCES decks(id) ON DELETE CASCADE
		);
		`,
		`
		CREATE TABLE IF NOT EXISTS card_tags(
			id INTEGER PRIMARY KEY,
			card_id INTEGER NOT NULL,
			tag_id INTEGER NOT NULL,
			FOREIGN KEY(card_id) REFERENCES cards(id) ON DELETE CASCADE,
			FOREIGN KEY(tag_id) REFERENCES tags(id) ON DELETE CASCADE
		);
		`,
	}

	for _, q := range queries {
		_, err = db.Exec(q)

		if err != nil {
			return nil, err
		}
	}

	return &Database{db}, nil
}

func (db *Database) Create(r web.Record) error {
	err := r.GenerateSlug()
	if err != nil {
		return errors.Wrapf(err, "failed to generate slug for %v", r)
	}

	rv := reflect.ValueOf(r)

	if rv.Kind() != reflect.Ptr {
		return errors.Errorf("cannot create database record for %v", r)
	}

	rv = rv.Elem()
	rt := reflect.TypeOf(rv.Interface())

	var columns []string
	var vars []string
	var fields []interface{}

	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		tag := f.Tag.Get("db")

		if tag == "" || tag == "id" {
			continue
		}

		columns = append(columns, tag)
		vars = append(vars, "?")
		fields = append(fields, rv.Field(i).Interface())
	}

	q := fmt.Sprintf("INSERT into %s (%s) values (%s);", r.Type(),
		strings.Join(columns, ", "), strings.Join(vars, ", "))

	res, err := db.Exec(q, fields...)

	if err != nil {
		return errors.Wrap(err, "failed to create db record")
	}

	id, err := res.LastInsertId()

	if err != nil {
		return errors.Wrap(err, "failed to retrieve last inserted id")
	}

	r.SetID(uint(id))

	return nil
}

func (db *Database) Query(where string, col web.Collection) error {
	r := col.NewRecord()
	rv := reflect.ValueOf(r)

	if rv.Kind() != reflect.Ptr {
		return errors.Errorf("cannot query database records for %v", col)
	}

	rv = rv.Elem()
	rt := reflect.TypeOf(rv.Interface())

	var columns []string

	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		tag := f.Tag.Get("db")

		if tag == "" {
			continue
		}

		columns = append(columns, tag)
	}

	if where != "" {
		where = "WHERE " + where
	}

	q := fmt.Sprintf("SELECT %s FROM %s %s;", strings.Join(columns, ", "), r.Type(), where)

	rows, err := db.DB.Query(q)
	if err != nil {
		return errors.Wrapf(err, "failed to query records %s", q)
	}
	defer rows.Close()

	for rows.Next() {

		r := col.NewRecord()
		rv := reflect.ValueOf(r)

		if rv.Kind() != reflect.Ptr {
			return errors.Errorf("cannot query database records for %v", col)
		}

		rv = rv.Elem()
		rt := reflect.TypeOf(rv.Interface())
		var fields []interface{}

		for i := 0; i < rt.NumField(); i++ {
			f := rt.Field(i)
			tag := f.Tag.Get("db")

			if tag == "" {
				continue
			}

			field := rv.Field(i).Addr().Interface()
			fields = append(fields, field)
		}

		err = rows.Scan(fields...)
		if err != nil {
			return errors.Wrap(err, "failed to scan records")
		}

		col.Append(r)
	}
	err = rows.Err()
	if err != nil {
		return errors.Wrap(err, "find to query records")
	}
	return nil
}

func (db *Database) Get(slug string, r web.Record) error {
	rv := reflect.ValueOf(r)

	if rv.Kind() != reflect.Ptr {
		return errors.Errorf("cannot query database records for %v", r)
	}

	rv = rv.Elem()
	rt := reflect.TypeOf(rv.Interface())

	var columns []string
	var fields []interface{}

	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		tag := f.Tag.Get("db")

		if tag == "" {
			continue
		}

		columns = append(columns, tag)
		field := rv.Field(i).Addr().Interface()
		fields = append(fields, field)
	}

	q := fmt.Sprintf("SELECT %s FROM %s WHERE slug = %q LIMIT 1;",
		strings.Join(columns, ", "), r.Type(), slug)

	rows, err := db.DB.Query(q)
	if err != nil {
		return errors.Wrapf(err, "failed to get record %s", q)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(fields...)
		if err != nil {
			return errors.Wrap(err, "failed to scan record")
		}

		return nil
	}
	err = rows.Err()
	if err != nil {
		return errors.Wrap(err, "find to get record")
	}
	return nil
}
