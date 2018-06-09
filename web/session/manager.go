package session

import (
	"net/http"

	"github.com/pkg/errors"
	"gitlab.com/luizbranco/srs/db"
	"gitlab.com/luizbranco/srs/primitives"
)

const cookieName = "session"

type Manager struct {
	Database primitives.Database
	Secret   string
}

func (m *Manager) LogIn(u primitives.User, w http.ResponseWriter) error {
	s := &session{
		UserID: u.ID(),
	}

	err := m.Database.Create(s)
	if err != nil {
		return errors.Wrap(err, "failed to create session")
	}

	value, err := m.encrypt(s.ID())
	if err != nil {
		return errors.Wrap(err, "failed to encrypt session id")
	}

	cookie := &http.Cookie{
		Name:  cookieName,
		Value: value,
		Path:  "/",
		//Secure:   true, // FIXME
		HttpOnly: true,
		MaxAge:   30 * 24 * 60 * 60, // 30 days
	}

	http.SetCookie(w, cookie)

	return nil
}

func (m *Manager) LogOut(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   cookieName,
		MaxAge: -1,
	}

	http.SetCookie(w, cookie)
}

func (m *Manager) User(r *http.Request) (*primitives.User, error) {
	var value string

	cookie, err := r.Cookie(cookieName)
	switch err {
	case nil: // found
		value = cookie.Value
	case http.ErrNoCookie:
		return nil, nil
	default:
		return nil, errors.Wrap(err, "invalid session cookie")
	}

	id, err := m.decrypt(value)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decrypt session id")
	}

	user, err := db.FindUser(m.Database, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find user")
	}

	return user, nil
}
