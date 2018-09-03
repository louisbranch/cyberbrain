package primitives

import (
	"time"
)

type CardReview struct {
	MetaID        ID        `db:"id"`
	MetaVersion   int       `db:"version"`
	MetaCreatedAt time.Time `db:"created_at"`
	MetaUpdatedAt time.Time `db:"updated_at"`

	DeckID  ID     `db:"deck_id"`
	CardID  ID     `db:"card_id"`
	Answer  string `db:"answer"`
	Skipped bool   `db:"skipped"`
	Correct bool   `db:"correct"`
}

func (c CardReview) ID() ID {
	return c.MetaID
}

func (c CardReview) Type() string {
	return "card_review"
}

func (c CardReview) Slug() string {
	return "review"
}

func (c *CardReview) SetID(id ID) {
	c.MetaID = id
}

func (c *CardReview) SetVersion(v int) {
	c.MetaVersion = v
}

func (c *CardReview) SetCreatedAt(t time.Time) {
	c.MetaCreatedAt = t
}

func (c *CardReview) SetUpdatedAt(t time.Time) {
	c.MetaUpdatedAt = t
}
