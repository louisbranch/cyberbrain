package primitives

import "time"

type User struct {
	MetaID        ID        `db:"id"`
	MetaVersion   int       `db:"version"`
	MetaCreatedAt time.Time `db:"created_at"`
	MetaUpdatedAt time.Time `db:"updated_at"`

	Email        string `db:"email"`
	Name         string `db:"name"`
	PasswordHash string `db:"password_hash"`
	ImageURL     string `db:"image_url"`
}

func (u User) ID() ID {
	return u.MetaID
}

func (u User) Type() string {
	return "user"
}

func (u *User) SetID(id ID) {
	u.MetaID = id
}

func (u *User) SetVersion(v int) {
	u.MetaVersion = v
}

func (u *User) SetCreatedAt(t time.Time) {
	u.MetaCreatedAt = t
}

func (u *User) SetUpdatedAt(t time.Time) {
	u.MetaUpdatedAt = t
}
