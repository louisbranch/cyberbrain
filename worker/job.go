package worker

import (
	"time"

	"gitlab.com/luizbranco/srs/primitives"
)

type Job struct {
	MetaID        primitives.ID `db:"id"`
	MetaVersion   int           `db:"version"`
	MetaCreatedAt time.Time     `db:"created_at"`
	MetaUpdatedAt time.Time     `db:"updated_at"`

	Name string `db:"name"`
	Args []byte `db:"args"`
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
