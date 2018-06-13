package session

import (
	"net/http"
	"strconv"

	"github.com/pkg/errors"
	"gitlab.com/luizbranco/cyberbrain/db"
	"gitlab.com/luizbranco/cyberbrain/primitives"
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
		Path:   "/",
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

	n, err := strconv.Atoi(id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read session id")
	}

	s, err := findSession(m.Database, primitives.ID(n))
	if err != nil {
		return nil, errors.Wrap(err, "failed to find session")
	}

	user, err := db.FindUser(m.Database, s.UserID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find user")
	}

	return user, nil
}

func findSession(db primitives.Database, id primitives.ID) (*session, error) {
	q := &query{id: id}

	r, err := db.Get(q)
	if err != nil {
		return nil, err
	}

	session, ok := r.(*session)
	if !ok {
		return nil, errors.Errorf("invalid record type %T", r)
	}

	return session, nil
}

type query struct {
	id primitives.ID
}

func (q *query) NewRecord() primitives.Record {
	return &session{}
}

func (q *query) Where() map[string]interface{} {
	return map[string]interface{}{
		"id": q.id,
	}
}

func (q *query) Raw() string {
	return ""
}

func (q *query) SortBy() map[string]string {
	return nil
}
