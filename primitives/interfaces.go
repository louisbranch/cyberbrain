package primitives

import "time"

type ID int

type Identifiable interface {
	ID() ID
	Type() string
}

type Database interface {
	Create(Record) error
	Update(Record) error
	Query(Query) ([]Record, error)
	QueryRaw(Query) ([]Record, error)
	Get(Query) (Record, error)
	Count(Query) (int, error)
	Random(Query, int) ([]Record, error)
}

type Record interface {
	Identifiable

	SetID(ID)
	SetVersion(int)
	SetCreatedAt(time.Time)
	SetUpdatedAt(time.Time)
}

type Query interface {
	NewRecord() Record
	Where() map[string]interface{}
	Raw() string
	SortBy() map[string]string
}

type PracticeGenerator interface {
	NewRound(Practice) (*Round, error)
}