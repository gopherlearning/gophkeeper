package client

import (
	"testing"

	"github.com/gopherlearning/gophkeeper/internal"
	"github.com/stretchr/testify/assert"
)

func Test_generateMnemonic(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, err := generateMnemonic()
		assert.NoError(t, err)
		assert.NotEmpty(t, m)
	})
	t.Run("NewEntropy error", func(t *testing.T) {
		entropyBitSize = 127
		m, err := generateMnemonic()
		assert.Error(t, err)
		assert.Empty(t, m)
		entropyBitSize = 128
	})
	t.Run("NewMnemonic error", func(t *testing.T) {
		internal.EmulatedError = "NewMnemonic error 1"
		m, err := generateMnemonic()
		assert.Error(t, err)
		assert.Empty(t, m)
		internal.EmulatedError = ""
	})
}
