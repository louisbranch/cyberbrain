package primitives

import (
	"time"
)

type Card struct {
	MetaID        ID        `db:"id"`
	MetaVersion   int       `db:"version"`
	MetaCreatedAt time.Time `db:"created_at"`
	MetaUpdatedAt time.Time `db:"updated_at"`

	DeckID      ID       `db:"deck_id"`
	Definitions []string `db:"definitions"`
	ImageURL    string   `db:"image_url"`
	SoundURL    string   `db:"sound_url"`
	Caption     string   `db:"caption"`
	NSFW        bool     `db:"nsfw"`
}

func (c Card) ID() ID {
	return c.MetaID
}

func (c Card) Type() string {
	return "card"
}

func (c *Card) SetID(id ID) {
	c.MetaID = id
}

func (c *Card) SetVersion(v int) {
	c.MetaVersion = v
}

func (c *Card) SetCreatedAt(t time.Time) {
	c.MetaCreatedAt = t
}

func (c *Card) SetUpdatedAt(t time.Time) {
	c.MetaUpdatedAt = t
}

func (c *Card) GetImageURL() string {
	return c.ImageURL
}

func (c *Card) SetImageURL(url string) {
	c.ImageURL = url
}
