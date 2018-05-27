package web

import "fmt"

const (
	PracticeStateInProgress = "in_progress"
)

type Practice struct {
	ID     uint   `db:"id"`
	DeckID uint   `db:"deck_id"`
	Rounds int    `db:"rounds"`
	State  string `db:"state"`
}

func (p *Practice) Type() string {
	return "practices"
}

func (p *Practice) SetID(id uint) {
	p.ID = id
}

func (p *Practice) GenerateSlug() error {
	return nil
}

func (p *Practice) Slug() string {
	return fmt.Sprintf("%d", p.ID)
}

type Practices []Practice

func (p *Practices) NewRecord() Record {
	return &Practice{}
}

func (p *Practices) Append(r Record) {
	practice := r.(*Practice)
	*p = append(*p, *practice)
}

type PracticeRound struct {
	ID         uint   `db:"id"`
	PracticeID uint   `db:"practice_id"`
	CardID     uint   `db:"card_id"`
	Round      uint   `db:"round"`
	Expect     string `db:"expect"`
	Answer     string `db:"answer"`
	Correct    bool   `db:"correct"`
}

func (pr *PracticeRound) Type() string {
	return "practice_rounds"
}

func (pr *PracticeRound) SetID(id uint) {
	pr.ID = id
}

func (pr *PracticeRound) GenerateSlug() error {
	return nil
}

type PracticeRounds []PracticeRound

func (pr *PracticeRounds) NewRecord() Record {
	return &PracticeRound{}
}

func (pr *PracticeRounds) Append(r Record) {
	round := r.(*PracticeRound)
	*pr = append(*pr, *round)
}
