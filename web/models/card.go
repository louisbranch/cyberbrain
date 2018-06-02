package models

import (
	"net/url"
	"time"

	"github.com/luizbranco/srs/web"
)

type Card struct {
	MetaID        web.ID    `db:"id"`
	MetaVersion   int       `db:"version"`
	MetaCreatedAt time.Time `db:"created_at"`
	MetaUpdatedAt time.Time `db:"updated_at"`

	DeckID      web.ID   `db:"deck_id"`
	ImageURLs   []string `db:"image_urls"`
	AudioURLs   []string `db:"audio_urls"`
	Definitions []string `db:"definitions"`
	Tags        []Tag
}

func NewCardFromForm(deckID web.ID, form url.Values) (*Card, error) {
	c := &Card{
		DeckID: deckID,
	}

	for _, f := range form["image_urls"] {
		if f != "" {
			c.ImageURLs = append(c.ImageURLs, f)
		}
	}

	for _, f := range form["audio_urls"] {
		if f != "" {
			c.AudioURLs = append(c.AudioURLs, f)
		}
	}

	for _, f := range form["definitions"] {
		if f != "" {
			c.Definitions = append(c.Definitions, f)
		}
	}

	return c, nil
}

func (c *Card) ID() web.ID {
	return c.MetaID
}

func (c *Card) SetID(id web.ID) {
	c.MetaID = id
}

func (c *Card) Type() string {
	return "card"
}
