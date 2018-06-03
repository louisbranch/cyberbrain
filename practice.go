package srs

import (
	"time"
)

type PracticeMode int

const (
	PracticeImages PracticeMode = 1 << iota
	PracticeSounds
	PracticeDefinitions
	PracticeTags

	PracticeRandom PracticeMode = 0
)

type Practice struct {
	MetaID        ID        `db:"id"`
	MetaVersion   int       `db:"version"`
	MetaCreatedAt time.Time `db:"created_at"`
	MetaUpdatedAt time.Time `db:"updated_at"`

	DeckID ID           `db:"deck_id"`
	Mode   PracticeMode `db:"mode"`
	Field  int          `db:"field"`
	TagID  *ID          `db:"tag_id"`
	Rounds int          `db:"rounds"`
	Done   bool         `db:"done"`

	Deck           *Deck
	PracticeRounds []PracticeRound
}

func (p Practice) ID() ID {
	return p.MetaID
}

func (p Practice) Type() string {
	return "practice"
}

func (p *Practice) SetID(id ID) {
	p.MetaID = id
}

func (p *Practice) SetVersion(v int) {
	p.MetaVersion = v
}

func (p *Practice) SetCreatedAt(t time.Time) {
	p.MetaCreatedAt = t
}

func (p *Practice) SetUpdatedAt(t time.Time) {
	p.MetaUpdatedAt = t
}

type PracticeRound struct {
	MetaID        ID        `db:"id"`
	MetaVersion   int       `db:"version"`
	MetaCreatedAt time.Time `db:"created_at"`
	MetaUpdatedAt time.Time `db:"updated_at"`

	PracticeID ID           `db:"practice_id"`
	Mode       PracticeMode `db:"mode"`
	CardIDs    []ID         `db:"card_ids"`
	Options    []string     `db:"options"`
	Answer     string       `db:"answer"`
	Correct    bool         `db:"correct"`
	Done       bool         `db:"done"`
}

func (pr PracticeRound) ID() ID {
	return pr.MetaID
}

func (pr PracticeRound) Type() string {
	return "practice_round"
}

func (pr *PracticeRound) SetID(id ID) {
	pr.MetaID = id
}

func (pr *PracticeRound) SetVersion(v int) {
	pr.MetaVersion = v
}

func (pr *PracticeRound) SetCreatedAt(t time.Time) {
	pr.MetaCreatedAt = t
}

func (pr *PracticeRound) SetUpdatedAt(t time.Time) {
	pr.MetaUpdatedAt = t
}
