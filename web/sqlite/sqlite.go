package sqlite

import (
	"database/sql"
	"fmt"
	"strconv"

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
			slug TEXT NOT NULL UNIQUE CHECK(slug <> ''),
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
			slug TEXT NOT NULL UNIQUE CHECK(slug <> ''),
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

	last, err := res.LastInsertId()
	if err != nil {
		return errors.Wrap(err, "failed to retrieve last inserted id")
	}

	id := strconv.FormatInt(last, 10)

	r.SetID(web.ID(id))

	return nil
}

func (db *Database) Query(wq web.Query) ([]web.Record, error) {
	q, err := QueryFromRecord(wq.NewRecord())
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get record fields %s", q)
	}

	raw := fmt.Sprintf("SELECT %s FROM %s %s;", q.Columns(), q.Table(), where(wq))

	return db.queryRows(wq, raw)
}

func (db *Database) Get(wq web.Query) (web.Record, error) {
	r := wq.NewRecord()

	q, err := QueryFromRecord(r)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get record fields %s", q)
	}

	query := fmt.Sprintf("SELECT %s FROM %s %s;", q.Columns(), q.Table(), where(wq))

	row := db.DB.QueryRow(query)

	err = row.Scan(q.fields...)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to scan record %q", query)
	}

	return r, nil
}

func (db *Database) QueryRaw(wq web.Query) ([]web.Record, error) {
	return db.queryRows(wq, wq.Raw())
}

func (db *Database) Count(wq web.Query) (int, error) {
	table := wq.NewRecord().Type() + "s"

	query := fmt.Sprintf("SELECT COUNT(*) FROM %s %s;", table, where(wq))

	row := db.DB.QueryRow(query)

	var n int

	err := row.Scan(&n)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to count records %q", query)
	}

	return n, nil
}

func (db *Database) Random(wq web.Query, n int) ([]web.Record, error) {
	r := wq.NewRecord()
	q, err := QueryFromRecord(r)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get record fields %v", wq)
	}

	raw := fmt.Sprintf(`SELECT %s FROM %s WHERE id IN (SELECT id FROM %s ORDER BY RANDOM() LIMIT %d)`,
		q.Columns(), q.Table(), q.Table(), n)

	return db.queryRows(wq, raw)
}

func (db *Database) queryRows(wq web.Query, query string) ([]web.Record, error) {
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to query records %v", wq)
	}
	defer rows.Close()

	var records []web.Record

	for rows.Next() {

		r := wq.NewRecord()

		q, err := QueryFromRecord(r)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get record fields %v", wq)
		}

		err = rows.Scan(q.fields...)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to scan records %v", wq)
		}

		records = append(records, r)
	}

	err = rows.Err()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to query records %v", wq)
	}

	return records, nil
}
