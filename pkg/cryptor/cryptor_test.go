package cryptor

import (
	"crypto/ed25519"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var tests = []struct {
	name       string
	secret     ed25519.PrivateKey
	recipients []ed25519.PublicKey
	err        error
	res        string
	key        []byte
}{
	{
		name:   "too long",
		secret: []byte("100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
		err:    errors.New("too long"),
		key:    []byte{},
	},
	{
		name:   "short key",
		secret: []byte("sh"),
		err:    errors.New("invalid X25519 secret key"),
		key:    []byte{},
	},
	{
		name:       "success",
		recipients: []ed25519.PublicKey{[]byte("")},
		secret:     []byte{195, 47, 118, 28, 63, 49, 142, 147, 228, 15, 201, 172, 47, 106, 168, 170, 225, 118, 163, 101, 52, 186, 153, 36, 234, 160, 195, 215, 117, 201, 165, 58},
		err:        nil,
		key:        []byte{0xc3, 0x2f, 0x76, 0x1c, 0x3f, 0x31, 0x8e, 0x93, 0xe4, 0xf, 0xc9, 0xac, 0x2f, 0x6a, 0xa8, 0xaa, 0xe1, 0x76, 0xa3, 0x65, 0x34, 0xba, 0x99, 0x24, 0xea, 0xa0, 0xc3, 0xd7, 0x75, 0xc9, 0xa5, 0x3a},
	},
}

func TestCryptor(t *testing.T) {
	for _, v := range tests {
		c, err := NewCryptor(v.secret)
		if v.err != nil {
			assert.ErrorContains(t, err, v.err.Error())
			assert.Nil(t, c)
		} else {
			fmt.Println(string(v.key))
			fmt.Println(string(c.GetKey()))
			assert.NotNil(t, c)
			assert.Equal(t, v.key, c.GetKey())
			chip, err := c.Encrypt([]byte("test"))
			assert.NoError(t, err)
			assert.Equal(t, chip, []byte("test"))
			plain, err := c.Decrypt([]byte("test"))
			assert.NoError(t, err)
			assert.Equal(t, plain, []byte("test"))
		}
	}
}
