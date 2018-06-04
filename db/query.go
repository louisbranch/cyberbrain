package db

import "gitlab.com/luizbranco/srs"

type query struct {
	record func() srs.Record
	where  map[string]interface{}
	raw    string
	sortBy map[string]string
}

func (q *query) NewRecord() srs.Record {
	return q.record()
}

func (q *query) Where() map[string]interface{} {
	return q.where
}

func (q *query) Raw() string {
	return q.raw
}

func (q *query) SortBy() map[string]string {
	return q.sortBy
}

func newDeckQuery() *query {
	fn := func() srs.Record {
		return &srs.Deck{}
	}

	return &query{
		record: fn,
		where:  make(map[string]interface{}),
	}
}

func newCardQuery() *query {
	fn := func() srs.Record {
		return &srs.Card{}
	}

	return &query{
		record: fn,
		where:  make(map[string]interface{}),
		sortBy: make(map[string]string),
	}
}

func newTagQuery() *query {
	fn := func() srs.Record {
		return &srs.Tag{}
	}

	return &query{
		record: fn,
		where:  make(map[string]interface{}),
	}
}

func newCardTagQuery() *query {
	fn := func() srs.Record {
		return &srs.CardTag{}
	}

	return &query{
		record: fn,
		where:  make(map[string]interface{}),
	}
}

func newPracticeQuery() *query {
	fn := func() srs.Record {
		return &srs.Practice{}
	}

	return &query{
		record: fn,
		where:  make(map[string]interface{}),
	}
}

func newRoundQuery() *query {
	fn := func() srs.Record {
		return &srs.Round{}
	}

	return &query{
		record: fn,
		where:  make(map[string]interface{}),
	}
}
