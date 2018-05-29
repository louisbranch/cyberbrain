package web

import "fmt"

const (
	PracticeStateInProgress = "in_progress"
	PracticeStateFinished   = "finished"
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

func (p *Practice) Finished() bool {
	return p.State == PracticeStateFinished
}

type PracticeRound struct {
	ID         uint   `db:"id"`
	PracticeID uint   `db:"practice_id"`
	CardID     uint   `db:"card_id"`
	Round      int    `db:"round"`
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
