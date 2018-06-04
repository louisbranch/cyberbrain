package generator

import (
	"math/rand"
	"time"

	"github.com/pkg/errors"
	"gitlab.com/luizbranco/srs"
	"gitlab.com/luizbranco/srs/db"
)

type Generator struct{}

var r *rand.Rand

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func (Generator) NewRound(conn srs.Database, p srs.Practice) (*srs.Round, error) {
	pid := p.ID()

	// TODO query should find by tags and where field(s) are not null
	card, err := db.RandomCard(conn, p.DeckID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find random card for practice %q", pid)
	}

	// TODO implement different types of prompt and guess modes
	r := &srs.Round{
		PracticeID: pid,
		PromptMode: p.PromptMode,
		GuessMode:  p.GuessMode,
		CardIDs:    []srs.ID{card.ID()},
		Prompt:     card.ImageURLs[0],
		Answer:     card.Definitions[0],
	}

	err = conn.Create(r)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create round for practice %q, %v", pid, r)
	}

	return r, nil
}
