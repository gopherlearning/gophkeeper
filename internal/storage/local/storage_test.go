package local

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewStorage(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		err  error
		name string
		key  string
		path string
	}{
		{
			name: "success",
			key:  "secretwqwwwwwwww",
			path: t.TempDir(),
			err:  nil,
		},
		{
			name: "short key",
			key:  "secretwqwww",
			path: t.TempDir(),
			err:  errors.New("Encryption key's length should beeither"),
		},
		{
			name: "long key",
			key:  "dwdwddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd",
			path: t.TempDir(),
			err:  errors.New("Encryption key's length should beeither"),
		},
		{
			name: "bad dir",
			key:  "secretwqwwwwwwww",
			path: "",
			err:  errors.New("no such file or directory"),
		},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			s, err := NewLocalStorage(v.key, v.path, nil)
			if v.err != nil {
				require.ErrorContains(t, err, v.err.Error())
				assert.Nil(t, s)
				return
			}
			require.NoError(t, err)
			assert.NotNil(t, s)
			assert.NoError(t, s.Close())
		})
	}
}
