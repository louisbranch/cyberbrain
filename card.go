package srs

import (
	"time"
)

type Card struct {
	MetaID        ID        `db:"id"`
	MetaVersion   int       `db:"version"`
	MetaCreatedAt time.Time `db:"created_at"`
	MetaUpdatedAt time.Time `db:"updated_at"`

	DeckID      ID       `db:"deck_id"`
	ImageURLs   []string `db:"image_urls"`
	AudioURLs   []string `db:"audio_urls"`
	Definitions []string `db:"definitions"`
	Tags        []Tag
}

func (c Card) ID() ID {
	return c.MetaID
}

func (c *Card) SetID(id ID) {
	c.MetaID = id
}

func (c Card) Type() string {
	return "card"
}
