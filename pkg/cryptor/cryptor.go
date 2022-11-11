package cryptor

import (
	"bytes"
	"crypto/ed25519"
	"errors"
	"io"
	"strings"

	"filippo.io/age"
	"github.com/gopherlearning/gophkeeper/internal/bech32"
	"github.com/rs/zerolog/log"
)

var (
	ErrRecipientsCount = errors.New("необходим минимум 1 получатель")
)

type Cryptor struct {
	identity *age.X25519Identity
	secret   []byte
}

func (c *Cryptor) GetKey() []byte {
	return c.secret
}
func (c *Cryptor) Decrypt(chipher []byte) ([]byte, error) {
	plain, err := age.Decrypt(bytes.NewReader(chipher), c.identity)
	if err != nil {
		return nil, err
	}

	res, err := io.ReadAll(plain)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Cryptor) Encrypt(plaindata []byte, recipients ...string) ([]byte, error) {
	if len(recipients) == 0 {
		return nil, ErrRecipientsCount
	}

	rr := make([]age.Recipient, len(recipients))

	for i, v := range recipients {
		r, err := age.ParseX25519Recipient(v)
		if err != nil {
			return nil, err
		}

		rr[i] = r
	}

	out := &bytes.Buffer{}

	w, err := age.Encrypt(out, rr...)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(w, bytes.NewReader(plaindata))
	if err != nil {
		return nil, err
	}

	err = w.Close()
	if err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}

func NewCryptor(secret ed25519.PrivateKey) (*Cryptor, error) {
	key, err := toAgePrivate(secret)
	if err != nil {
		log.Debug().Err(err)
		return nil, err
	}

	identity, err := age.ParseX25519Identity(key)
	if err != nil {
		log.Debug().Err(err)
		return nil, err
	}

	return &Cryptor{identity: identity, secret: secret}, nil
}

func toAgePrivate(priv ed25519.PrivateKey) (string, error) {
	s, err := bech32.Encode("AGE-SECRET-KEY-", priv)
	if err != nil {
		log.Debug().Err(err)
		return "", err
	}

	return strings.ToUpper(s), nil
}

// func toAgePublic(pub ed25519.PublicKey) (string, error) {
// 	s, err := bech32.Encode("age", pub)
// 	if err != nil {
// 		return "", err
// 	}

// 	return s, nil
// }
