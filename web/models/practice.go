package models

import (
	"net/url"
	"strconv"

	"github.com/luizbranco/srs/web"
	"github.com/pkg/errors"
)

const (
	PracticeStateInProgress = "in_progress"
	PracticeStateFinished   = "finished"
)

type Practice struct {
	ID     web.ID `db:"id"`
	DeckID web.ID `db:"deck_id"`
	Slug   string `db:"slug"`
	Rounds int    `db:"rounds"`
	State  string `db:"state"`
}

func NewPractice() *Practice {
	return &Practice{Slug: NewSlug()}
}

func NewPracticeFromForm(deckID web.ID, form url.Values) (*Practice, error) {
	rounds := form.Get("rounds")
	n, err := strconv.Atoi(rounds)
	if err != nil {
		return nil, errors.Wrap(err, "invalid number of rounds")
	}

	p := NewPractice()
	p.DeckID = deckID
	p.Rounds = n
	p.State = PracticeStateInProgress

	return p, nil
}

func (p *Practice) SetID(id web.ID) {
	p.ID = id
}

func (p *Practice) Type() string {
	return "practice"
}

func (p *Practice) Finished() bool {
	return p.State == PracticeStateFinished
}

type PracticeRound struct {
	ID         web.ID `db:"id"`
	PracticeID web.ID `db:"practice_id"`
	CardID     web.ID `db:"card_id"`
	Round      int    `db:"round"`
	Expect     string `db:"expect"`
	Answer     string `db:"answer"`
	Correct    bool   `db:"correct"`
}

func (pr *PracticeRound) SetID(id web.ID) {
	pr.ID = id
}

func (pr *PracticeRound) Type() string {
	return "practice_round"
}
