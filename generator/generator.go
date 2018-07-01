package generator

import (
	"math/rand"
	"time"

	"github.com/pkg/errors"
	"gitlab.com/luizbranco/cyberbrain/db"
	"gitlab.com/luizbranco/cyberbrain/primitives"
)

type Generator struct {
	Database primitives.Database
}

var r *rand.Rand

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func (g Generator) NewRound(p primitives.Practice) (*primitives.Round, error) {
	pid := p.ID()

	// TODO query should find by tags and where field(s) are not null
	// TODO should exclude cards already used in the same practice
	card, err := db.RandomCard(g.Database, p.DeckID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find random card for practice %q", pid)
	}

	// TODO implement different types of prompt and guess modes
	r := &primitives.Round{
		PracticeID: pid,
		PromptMode: p.PromptMode,
		Caption:    card.Caption,
		GuessMode:  p.GuessMode,
		CardIDs:    []primitives.ID{card.ID()},
		Prompt:     card.ImageURL,
		Answer:     card.Definitions[0],
	}

	err = g.Database.Create(r)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create round for practice %d, %v", pid, r)
	}

	return r, nil
}
