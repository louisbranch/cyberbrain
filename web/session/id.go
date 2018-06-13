// adapted from github.com/gtank/cryptopasta
package session

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"
	"strconv"

	"github.com/pkg/errors"
	"gitlab.com/luizbranco/cyberbrain/primitives"
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

	return hex.EncodeToString(cypher), nil
}

func (m *Manager) decrypt(cipherID string) (id string, err error) {
	key := []byte(m.Secret)

	block, err := aes.NewCipher(key[:])
	if err != nil {
		return id, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return id, errors.Wrap(err, "invalid id")
	}

	decoded, err := hex.DecodeString(cipherID)
	if err != nil {
		return id, errors.Wrap(err, "failed to decode session id")
	}

	if len(decoded) < gcm.NonceSize() {
		return id, errors.New("malformed id")
	}

	b, err := gcm.Open(nil, decoded[:gcm.NonceSize()], decoded[gcm.NonceSize():], nil)
	if err != nil {
		return id, errors.Wrap(err, "tampered id")
	}

	return string(b), nil
}
