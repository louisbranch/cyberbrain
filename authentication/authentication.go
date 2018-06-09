package authentication

import (
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type Authenticator struct{}

func (a Authenticator) Create(password string) (string, error) {
	p := []byte(password)

	h, err := bcrypt.GenerateFromPassword(p, bcrypt.MinCost)
	if err != nil {
		return "", errors.Wrap(err, "failed to hash password")
	}

	return string(h), nil
}

func (a Authenticator) Verify(hash, password string) (bool, error) {
	h := []byte(hash)
	p := []byte(password)

	err := bcrypt.CompareHashAndPassword(h, p)
	if err != nil {
		return false, errors.Wrap(err, "failed to verify password")
	}

	return true, nil
}
