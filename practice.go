package srs

import (
	"time"
)

const (
	PracticeStateInProgress = "in_progress"
	PracticeStateFinished   = "finished"
)

type Practice struct {
	MetaID        ID        `db:"id"`
	MetaVersion   int       `db:"version"`
	MetaCreatedAt time.Time `db:"created_at"`
	MetaUpdatedAt time.Time `db:"updated_at"`

	DeckID ID     `db:"deck_id"`
	Mode   string `db:"mode"`
	Rounds int    `db:"rounds"`
	State  string `db:"state"`
}

func (p Practice) ID() ID {
	return p.MetaID
}

func (p *Practice) SetID(id ID) {
	p.MetaID = id
}

func (p Practice) Type() string {
	return "practice"
}

func (p Practice) Finished() bool {
	return p.State == PracticeStateFinished
}

type PracticeRound struct {
	MetaID        ID        `db:"id"`
	MetaVersion   int       `db:"version"`
	MetaCreatedAt time.Time `db:"created_at"`
	MetaUpdatedAt time.Time `db:"updated_at"`

	PracticeID ID       `db:"practice_id"`
	Mode       string   `db:"mode"`
	CardIDs    []ID     `db:"card_ids"`
	Options    []string `db:"options"`
	Answer     string   `db:"answer"`
	Correct    bool     `db:"correct"`
}

func (pr PracticeRound) ID() ID {
	return pr.MetaID
}

func (pr *PracticeRound) SetID(id ID) {
	pr.MetaID = id
}

func (pr PracticeRound) Type() string {
	return "practice_round"
}
