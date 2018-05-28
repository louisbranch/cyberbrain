package sqlite

import (
	"database/sql"
	"fmt"

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
			image_url TEXT
		);
		`,
		`
		CREATE TABLE IF NOT EXISTS cards(
			id INTEGER PRIMARY KEY,
			deck_id INTEGER NOT NULL,
			image_url TEXT NOT NULL CHECK(image_url <> ''),
			audio_url TEXT,
			definition TEXT NOT NULL CHECK(definition <> ''),
			alt_definition TEXT,
			pronunciation TEXT,
			FOREIGN KEY(deck_id) REFERENCES decks(id) ON DELETE CASCADE
		);
		`,
		`
		CREATE TABLE IF NOT EXISTS tags(
			id INTEGER PRIMARY KEY,
			deck_id INTEGER NOT NULL,
			slug TEXT NOT NULL UNIQUE CHECK(slug <> ''),
			name TEXT NOT NULL UNIQUE CHECK(name <> ''),
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
		`
		CREATE TABLE IF NOT EXISTS practices(
			id INTEGER PRIMARY KEY,
			deck_id INTEGER NOT NULL,
			state TEXT NOT NULL CHECK(state <> ''),
			rounds INTEGER NOT NULL CHECK(rounds > 0),
			FOREIGN KEY(deck_id) REFERENCES decks(id) ON DELETE CASCADE
		);
		`,
		`
		CREATE TABLE IF NOT EXISTS practice_rounds(
			id INTEGER PRIMARY KEY,
			practice_id INTEGER NOT NULL,
			card_id INTEGER NOT NULL,
			round INTEGER NOT NULL CHECK(round > 0),
			expect TEXT NOT NULL CHECK(expect <> ''),
			answer TEXT,
			correct BOOLEAN,
			FOREIGN KEY(practice_id) REFERENCES practices(id) ON DELETE CASCADE,
			FOREIGN KEY(card_id) REFERENCES cards(id) ON DELETE CASCADE
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

	q, err := QueryFromRecord(r, "id")
	if err != nil {
		return errors.Wrapf(err, "failed to get record fields %v", r)
	}

	query := fmt.Sprintf("INSERT into %s (%s) values (%s);", q.Table(), q.Columns(),
		q.Placeholders())

	res, err := db.Exec(query, q.fields...)

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

func (db *Database) Query(w web.Where, rs web.Records) error {
	r := rs.NewRecord()
	q, err := QueryFromRecord(r)
	if err != nil {
		return errors.Wrapf(err, "failed to get record fields %s", q)
	}

	query := fmt.Sprintf("SELECT %s FROM %s %s;", q.Columns(), q.Table(), where(w))

	return db.QueryRaw(query, rs)
}

func (db *Database) Get(w web.Where, r web.Record) error {
	q, err := QueryFromRecord(r)
	if err != nil {
		return errors.Wrapf(err, "failed to get record fields %s", q)
	}

	query := fmt.Sprintf("SELECT %s FROM %s %s;", q.Columns(), q.Table(), where(w))

	row := db.DB.QueryRow(query)

	err = row.Scan(q.fields...)
	if err != nil {
		return errors.Wrapf(err, "failed to scan record %q", query)
	}

	return nil
}

func (db *Database) QueryRaw(query string, rs web.Records) error {
	rows, err := db.DB.Query(query)
	if err != nil {
		return errors.Wrapf(err, "failed to query records %q", query)
	}
	defer rows.Close()

	for rows.Next() {

		r := rs.NewRecord()

		q, err := QueryFromRecord(r)
		if err != nil {
			return errors.Wrapf(err, "failed to get record fields %q", query)
		}

		err = rows.Scan(q.fields...)
		if err != nil {
			return errors.Wrapf(err, "failed to scan records %q", query)
		}

		rs.Append(r)
	}

	err = rows.Err()
	if err != nil {
		return errors.Wrapf(err, "failed to query records %q", query)
	}

	return nil
}

func (db *Database) Count(w web.Where, r web.Record) (int, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s %s;", r.Type(), where(w))

	row := db.DB.QueryRow(query)

	var n int

	err := row.Scan(&n)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to count records %q", query)
	}

	return n, nil
}

func (db *Database) Random(n int, rs web.Records) error {
	r := rs.NewRecord()
	q, err := QueryFromRecord(r)
	if err != nil {
		return errors.Wrapf(err, "failed to get record fields %v", r)
	}

	query := fmt.Sprintf(`SELECT %s FROM %s WHERE id IN (SELECT id FROM %s ORDER BY RANDOM() LIMIT %d)`,
		q.Columns(), q.Table(), q.Table(), n)

	return db.QueryRaw(query, rs)
}
