package worker

import (
	"time"

	"gitlab.com/luizbranco/srs/primitives"
)

const (
	scheduled = "scheduled"
	running   = "running"
	done      = "done"
	retry     = "retry"
	failed    = "failed"
)

type Job struct {
	MetaID        primitives.ID `db:"id"`
	MetaVersion   int           `db:"version"`
	MetaCreatedAt time.Time     `db:"created_at"`
	MetaUpdatedAt time.Time     `db:"updated_at"`

	RunAt time.Time `db:"run_at"`
	Name  string    `db:"name"`
	State string    `db:"state"`
	Args  []byte    `db:"args"`
	Error string    `db:"error"`
	Tries int       `db:"tries"`
}

func (j Job) ID() primitives.ID {
	return j.MetaID
}

func (j Job) Type() string {
	return "job"
}

func (j *Job) SetID(id primitives.ID) {
	j.MetaID = id
}

func (j *Job) SetVersion(v int) {
	j.MetaVersion = v
}

func (j *Job) SetCreatedAt(t time.Time) {
	j.MetaCreatedAt = t
}

func (j *Job) SetUpdatedAt(t time.Time) {
	j.MetaUpdatedAt = t
}
