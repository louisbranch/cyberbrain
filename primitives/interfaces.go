package primitives

import (
	"context"
	"time"
)

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

type Authenticator interface {
	Create(password string) (string, error)
	Verify(hash string, password string) (bool, error)
}

type WorkerPool interface {
	Register(string, Worker) error
	Enqueue(string, map[string]string) error
}

type Worker interface {
	Spawn(map[string]string) (Job, error)
}

type Job interface {
	Run(ctx context.Context) error
}
