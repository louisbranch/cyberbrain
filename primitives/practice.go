package primitives

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

	DeckID      ID           `db:"deck_id"`
	PromptMode  PracticeMode `db:"prompt_mode"`
	GuessMode   PracticeMode `db:"guess_mode"`
	Field       int          `db:"field"`
	TagID       *ID          `db:"tag_id"`
	TotalRounds int          `db:"total_rounds"`
	Done        bool         `db:"done"`

	Deck   *Deck
	Rounds []Round
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
