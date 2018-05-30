package psql

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/luizbranco/srs/web"
	"github.com/pkg/errors"
)

type Database struct {
	*sql.DB
}

func New(host, port, dbname, user, pass string) (*Database, error) {
	params := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable", host,
		port, dbname, user, pass)

	db, err := sql.Open("postgres", params)
	if err != nil {
		return nil, err
	}

	queries := []string{
		`
		CREATE TABLE IF NOT EXISTS decks(
			id SERIAL PRIMARY KEY,
			slug TEXT NOT NULL UNIQUE CHECK(slug <> ''),
			name TEXT NOT NULL CHECK(name <> ''),
			description TEXT,
			image_url TEXT
		);
		`,
		`
		CREATE TABLE IF NOT EXISTS cards(
			id SERIAL PRIMARY KEY,
			deck_id INTEGER NOT NULL REFERENCES decks ON DELETE CASCADE,
			slug TEXT NOT NULL UNIQUE CHECK(slug <> ''),
			image_url TEXT NOT NULL CHECK(image_url <> ''),
			audio_url TEXT,
			definition TEXT NOT NULL CHECK(definition <> ''),
			alt_definition TEXT,
			pronunciation TEXT
		);
		`,
		`
		CREATE TABLE IF NOT EXISTS tags(
			id SERIAL PRIMARY KEY,
			deck_id INTEGER NOT NULL REFERENCES decks ON DELETE CASCADE,
			slug TEXT NOT NULL UNIQUE CHECK(slug <> ''),
			name TEXT NOT NULL UNIQUE CHECK(name <> '')
		);
		`,
		`
		CREATE TABLE IF NOT EXISTS card_tags(
			id SERIAL PRIMARY KEY,
			card_id INTEGER NOT NULL REFERENCES cards ON DELETE CASCADE,
			tag_id INTEGER NOT NULL REFERENCES tags ON DELETE CASCADE,
			UNIQUE (card_id, tag_id)
		);
		`,
		`
		CREATE TABLE IF NOT EXISTS practices(
			id SERIAL PRIMARY KEY,
			deck_id INTEGER NOT NULL REFERENCES decks ON DELETE CASCADE,
			slug TEXT NOT NULL UNIQUE CHECK(slug <> ''),
			state TEXT NOT NULL CHECK(state <> ''),
			rounds INTEGER NOT NULL CHECK(rounds > 0)
		);
		`,
		`
		CREATE TABLE IF NOT EXISTS practice_rounds(
			id SERIAL PRIMARY KEY,
			card_id INTEGER NOT NULL REFERENCES cards ON DELETE CASCADE,
			practice_id INTEGER NOT NULL REFERENCES practices ON DELETE CASCADE,
			round INTEGER NOT NULL CHECK(round > 0),
			expect TEXT NOT NULL CHECK(expect <> ''),
			answer TEXT,
			correct BOOLEAN,
			UNIQUE (card_id, practice_id)
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

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURNING id;", q.Table(), q.Columns(),
		q.Placeholders())

	var id web.ID

	err = db.QueryRow(query, q.fields...).Scan(&id)
	if err != nil {
		return errors.Wrapf(err, "failed to create db record %q", query)
	}

	r.SetID(id)

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
