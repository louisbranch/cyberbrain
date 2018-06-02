package psql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/luizbranco/srs"
	"github.com/pkg/errors"
)

const version = 1

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

	for _, q := range tableQueries {
		_, err = db.Exec(q)

		if err != nil {
			return nil, err
		}
	}

	return &Database{db}, nil
}

func (db *Database) Create(r srs.Record) error {
	now := time.Now()

	r.SetCreatedAt(now)
	r.SetUpdatedAt(now)
	r.SetVersion(version)

	q, err := QueryFromRecord(r, "id")
	if err != nil {
		return errors.Wrapf(err, "failed to get record fields %v", r)
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURNING id;", q.Table(), q.Columns(),
		q.Placeholders())

	var id srs.ID

	err = db.QueryRow(query, q.addrs...).Scan(&id)
	if err != nil {
		return errors.Wrapf(err, "failed to create db record %q", query)
	}

	r.SetID(id)

	return nil
}

func (db *Database) Query(wq srs.Query) ([]srs.Record, error) {
	q, err := QueryFromRecord(wq.NewRecord())
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get record fields %s", q)
	}

	raw := fmt.Sprintf("SELECT %s FROM %s %s;", q.Columns(), q.Table(), where(wq))

	return db.queryRows(wq, raw)
}

func (db *Database) Get(wq srs.Query) (srs.Record, error) {
	r := wq.NewRecord()

	q, err := QueryFromRecord(r)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get record fields %s", q)
	}

	query := fmt.Sprintf("SELECT %s FROM %s %s;", q.Columns(), q.Table(), where(wq))

	row := db.DB.QueryRow(query)

	err = q.Scan(row)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to scan record %q", query)
	}

	return r, nil
}

func (db *Database) QueryRaw(wq srs.Query) ([]srs.Record, error) {
	return db.queryRows(wq, wq.Raw())
}

func (db *Database) Count(wq srs.Query) (int, error) {
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

func (db *Database) Random(wq srs.Query, n int) ([]srs.Record, error) {
	r := wq.NewRecord()
	q, err := QueryFromRecord(r)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get record fields %v", wq)
	}

	raw := fmt.Sprintf(`SELECT %s FROM %s WHERE id IN (SELECT id FROM %s ORDER BY RANDOM() LIMIT %d)`,
		q.Columns(), q.Table(), q.Table(), n)

	return db.queryRows(wq, raw)
}

func (db *Database) queryRows(wq srs.Query, query string) ([]srs.Record, error) {
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to query records %v", wq)
	}
	defer rows.Close()

	var records []srs.Record

	for rows.Next() {

		r := wq.NewRecord()

		q, err := QueryFromRecord(r)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get record fields %v", wq)
		}

		err = q.Scan(rows)
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
