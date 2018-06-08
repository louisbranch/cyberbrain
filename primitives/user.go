package primitives

import "time"

type User struct {
	MetaID        ID        `db:"id"`
	MetaVersion   int       `db:"version"`
	MetaCreatedAt time.Time `db:"created_at"`
	MetaUpdatedAt time.Time `db:"updated_at"`

	Email        string `db:"string"`
	Name         string `db:"string"`
	PasswordHash string `db:"password_hash"`
	ImageURL     string `db:"image_url"`
}
