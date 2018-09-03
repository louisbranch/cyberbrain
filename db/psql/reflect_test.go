package psql

import (
	"testing"
	"time"

	"gitlab.com/luizbranco/cyberbrain/db"
	"gitlab.com/luizbranco/cyberbrain/primitives"
)

type query struct {
	record primitives.Record
	where  map[string]interface{}
	raw    string
	sortBy map[string]string
}

func (q query) NewRecord() primitives.Record {
	return q.record
}

func (q query) Where() map[string]interface{} {
	return q.where
}

func (q query) Raw() string {
	return q.raw
}

func (q query) SortBy() map[string]string {
	return q.sortBy
}

func TestWhere(t *testing.T) {

	tcs := []struct {
		scenario string
		cond     primitives.Query
		clause   string
	}{
		{
			scenario: "date greater or equal",
			cond: query{
				where: map[string]interface{}{
					"next_date": db.GreaterOrEqual{Time: time.Date(2018, time.September, 2, 12, 0, 0, 0, time.UTC)},
				},
			},
			clause: "WHERE next_date >= '2018-09-02T12:00:00Z'::date",
		},
		{
			scenario: "date less or equal",
			cond: query{
				where: map[string]interface{}{
					"next_date": db.LessOrEqual{Time: time.Date(2018, time.September, 2, 12, 0, 0, 0, time.UTC)},
				},
			},
			clause: "WHERE next_date <= '2018-09-02T12:00:00Z'::date",
		},
		{
			scenario: "sort by asc",
			cond: query{
				sortBy: map[string]string{
					"created_at": "ASC",
				},
			},
			clause: " ORDER BY created_at ASC",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.scenario, func(t *testing.T) {
			w := where(tc.cond)

			if w != tc.clause {
				t.Errorf("expected where to eq %q, got %q", tc.clause, w)
			}
		})
	}

}
