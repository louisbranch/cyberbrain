package primitives

import (
	"math"
	"time"
)

type CardSchedule struct {
	MetaID        ID        `db:"id"`
	MetaVersion   int       `db:"version"`
	MetaCreatedAt time.Time `db:"created_at"`
	MetaUpdatedAt time.Time `db:"updated_at"`

	NextDate     time.Time `db:"next_date"`
	DeckID       ID        `db:"deck_id"`
	CardID       ID        `db:"card_id"`
	CurrentScore int       `db:"current_score"`
	MaxScore     int       `db:"max_score"`
}

func NewCardSchedule(deckID, cardID ID) *CardSchedule {
	return &CardSchedule{
		NextDate: days(1),
		DeckID:   deckID,
		CardID:   cardID,
	}
}

func (c CardSchedule) ID() ID {
	return c.MetaID
}

func (c CardSchedule) Type() string {
	return "card_schedule"
}

func (c *CardSchedule) SetID(id ID) {
	c.MetaID = id
}

func (c *CardSchedule) SetVersion(v int) {
	c.MetaVersion = v
}

func (c *CardSchedule) SetCreatedAt(t time.Time) {
	c.MetaCreatedAt = t
}

func (c *CardSchedule) SetUpdatedAt(t time.Time) {
	c.MetaUpdatedAt = t
}

func (c *CardSchedule) Reschedule(correct bool) {
	switch {
	case correct && c.CurrentScore >= 0:
		c.CurrentScore += 1
		c.MaxScore += 1
	case correct:
		c.CurrentScore = 1
	case c.CurrentScore >= 0:
		c.CurrentScore = -1
	default:
		c.CurrentScore -= 1
	}

	switch {
	case c.CurrentScore <= 0:
		c.NextDate = days(1)
	case c.CurrentScore == 1:
		c.NextDate = days(2)
	default:
		n := int(math.Pow(float64(c.CurrentScore), 2))
		c.NextDate = days(n)
	}
}

func days(d int) time.Time {
	now := time.Now().UTC()
	begin := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	return begin.AddDate(0, 0, d)
}
