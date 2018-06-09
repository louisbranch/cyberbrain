// adapted from github.com/gtank/cryptopasta
package session

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"strconv"

	"github.com/pkg/errors"
	"gitlab.com/luizbranco/srs/primitives"
)

func (m *Manager) encrypt(id primitives.ID) (cypherID string, err error) {
	key := []byte(m.Secret)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return "", err
	}

	i := strconv.Itoa(int(id))

	cypher := gcm.Seal(nonce, nonce, []byte(i), nil)

	return fmt.Sprintf("%x", cypher), nil
}

func (m *Manager) decrypt(cipherID string) (id primitives.ID, err error) {
	key := []byte(m.Secret)

	block, err := aes.NewCipher(key[:])
	if err != nil {
		return id, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return id, errors.Wrap(err, "invalid id")
	}

	cid := []byte(string(cipherID))

	if len(cid) < gcm.NonceSize() {
		return id, errors.New("malformed id")
	}

	b, err := gcm.Open(nil, cid[:gcm.NonceSize()], cid[gcm.NonceSize():], nil)
	if err != nil {
		return id, errors.Wrap(err, "tampered id")
	}

	n, err := strconv.Atoi(string(b))
	if err != nil {
		return id, errors.Wrap(err, "invalid id number")
	}

	return primitives.ID(n), nil
}
