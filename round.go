package srs

import "time"

type Round struct {
	MetaID        ID        `db:"id"`
	MetaVersion   int       `db:"version"`
	MetaCreatedAt time.Time `db:"created_at"`
	MetaUpdatedAt time.Time `db:"updated_at"`

	PracticeID ID           `db:"practice_id"`
	Mode       PracticeMode `db:"mode"`
	CardIDs    []ID         `db:"card_ids"`
	Options    []string     `db:"options"`
	Guess      string       `db:"guess"`
	Correct    bool         `db:"correct"`
	Done       bool         `db:"done"`
}

func (r Round) ID() ID {
	return r.MetaID
}

func (r Round) Type() string {
	return "round"
}

func (r *Round) SetID(id ID) {
	r.MetaID = id
}

func (r *Round) SetVersion(v int) {
	r.MetaVersion = v
}

func (r *Round) SetCreatedAt(t time.Time) {
	r.MetaCreatedAt = t
}

func (r *Round) SetUpdatedAt(t time.Time) {
	r.MetaUpdatedAt = t
}
