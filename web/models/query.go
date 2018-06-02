package models

import "github.com/luizbranco/srs/web"

type query struct {
	record func() web.Record
	where  map[string]interface{}
	raw    string
}

func (q *query) NewRecord() web.Record {
	return q.record()
}

func (q *query) Where() map[string]interface{} {
	return q.where
}

func (q *query) Raw() string {
	return q.raw
}

func newDeckQuery() *query {
	fn := func() web.Record {
		return &Deck{}
	}

	return &query{
		record: fn,
		where:  make(map[string]interface{}),
	}
}

func newCardQuery() *query {
	fn := func() web.Record {
		return &Card{}
	}

	return &query{
		record: fn,
		where:  make(map[string]interface{}),
	}
}

func newTagQuery() *query {
	fn := func() web.Record {
		return &Tag{}
	}

	return &query{
		record: fn,
		where:  make(map[string]interface{}),
	}
}

func newCardTagQuery() *query {
	fn := func() web.Record {
		return &CardTag{}
	}

	return &query{
		record: fn,
		where:  make(map[string]interface{}),
	}
}

func newPracticeQuery() *query {
	fn := func() web.Record {
		return &Practice{}
	}

	return &query{
		record: fn,
		where:  make(map[string]interface{}),
	}
}

func newPracticeRoundQuery() *query {
	fn := func() web.Record {
		return &PracticeRound{}
	}

	return &query{
		record: fn,
		where:  make(map[string]interface{}),
	}
}
