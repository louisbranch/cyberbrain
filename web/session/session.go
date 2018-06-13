package session

import (
	"time"

	"gitlab.com/luizbranco/cyberbrain/primitives"
)

type session struct {
	MetaID        primitives.ID `db:"id"`
	MetaVersion   int           `db:"version"`
	MetaCreatedAt time.Time     `db:"created_at"`
	MetaUpdatedAt time.Time     `db:"updated_at"`

	UserID primitives.ID `db:"user_id"`
}

func (s session) ID() primitives.ID {
	return s.MetaID
}

func (s session) Type() string {
	return "session"
}

func (s *session) SetID(id primitives.ID) {
	s.MetaID = id
}

func (s *session) SetVersion(v int) {
	s.MetaVersion = v
}

func (s *session) SetCreatedAt(t time.Time) {
	s.MetaCreatedAt = t
}

func (s *session) SetUpdatedAt(t time.Time) {
	s.MetaUpdatedAt = t
}
