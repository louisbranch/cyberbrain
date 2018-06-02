package models

import (
	"net/url"
	"strconv"
	"time"

	"github.com/luizbranco/srs/web"
	"github.com/pkg/errors"
)

const (
	PracticeStateInProgress = "in_progress"
	PracticeStateFinished   = "finished"
)

type Practice struct {
	MetaID        web.ID    `db:"id"`
	MetaVersion   int       `db:"version"`
	MetaCreatedAt time.Time `db:"created_at"`
	MetaUpdatedAt time.Time `db:"updated_at"`

	DeckID web.ID `db:"deck_id"`
	Rounds int    `db:"rounds"`
	State  string `db:"state"`
}

func NewPracticeFromForm(deckID web.ID, form url.Values) (*Practice, error) {
	rounds := form.Get("rounds")
	n, err := strconv.Atoi(rounds)
	if err != nil {
		return nil, errors.Wrap(err, "invalid number of rounds")
	}

	p := &Practice{
		DeckID: deckID,
		Rounds: n,
		State:  PracticeStateInProgress,
	}

	return p, nil
}

func (p *Practice) ID() web.ID {
	return p.MetaID
}

func (p *Practice) SetID(id web.ID) {
	p.MetaID = id
}

func (p *Practice) Type() string {
	return "practice"
}

func (p *Practice) Finished() bool {
	return p.State == PracticeStateFinished
}

type PracticeRound struct {
	MetaID     web.ID `db:"id"`
	PracticeID web.ID `db:"practice_id"`
	CardID     web.ID `db:"card_id"`
	Round      int    `db:"round"`
	Expect     string `db:"expect"`
	Answer     string `db:"answer"`
	Correct    bool   `db:"correct"`
}

func (pr *PracticeRound) ID() web.ID {
	return pr.MetaID
}

func (pr *PracticeRound) SetID(id web.ID) {
	pr.MetaID = id
}

func (pr *PracticeRound) Type() string {
	return "practice_round"
}
