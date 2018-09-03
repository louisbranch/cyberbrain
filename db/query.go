package db

import "gitlab.com/luizbranco/cyberbrain/primitives"

type query struct {
	record func() primitives.Record
	where  map[string]interface{}
	raw    string
	sortBy map[string]string
}

func (q *query) NewRecord() primitives.Record {
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

func newUserQuery() *query {
	fn := func() primitives.Record {
		return &primitives.User{}
	}

	return &query{
		record: fn,
		where:  make(map[string]interface{}),
		sortBy: make(map[string]string),
	}
}

func newDeckQuery() *query {
	fn := func() primitives.Record {
		return &primitives.Deck{}
	}

	return &query{
		record: fn,
		where:  make(map[string]interface{}),
		sortBy: make(map[string]string),
	}
}

func newCardQuery() *query {
	fn := func() primitives.Record {
		return &primitives.Card{}
	}

	return &query{
		record: fn,
		where:  make(map[string]interface{}),
		sortBy: make(map[string]string),
	}
}

func newTagQuery() *query {
	fn := func() primitives.Record {
		return &primitives.Tag{}
	}

	return &query{
		record: fn,
		where:  make(map[string]interface{}),
		sortBy: make(map[string]string),
	}
}

func newCardTagQuery() *query {
	fn := func() primitives.Record {
		return &primitives.CardTag{}
	}

	return &query{
		record: fn,
		where:  make(map[string]interface{}),
		sortBy: make(map[string]string),
	}
}

func newCardScheduleQuery() *query {
	fn := func() primitives.Record {
		return &primitives.CardSchedule{}
	}

	return &query{
		record: fn,
		where:  make(map[string]interface{}),
		sortBy: make(map[string]string),
	}
}

func newCardReviewQuery() *query {
	fn := func() primitives.Record {
		return &primitives.CardReview{}
	}

	return &query{
		record: fn,
		where:  make(map[string]interface{}),
		sortBy: make(map[string]string),
	}
}
