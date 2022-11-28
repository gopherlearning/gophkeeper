package local

import (
	"errors"
	"testing"

	"github.com/gopherlearning/gophkeeper/internal/model"
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
			s, err := NewLocalStorage(v.key, v.path, "")
			if v.err != nil {
				require.ErrorContains(t, err, v.err.Error())
				assert.Nil(t, s)
				return
			}
			status, ok := s.Status()()
			assert.Empty(t, status)
			assert.Equal(t, ok, false)
			s.remoteStatus.Store("123")
			status, ok = s.Status()()
			assert.Equal(t, "123", status)
			assert.Equal(t, ok, true)
			require.NoError(t, err)
			assert.NotNil(t, s)
			assert.NoError(t, s.Update(model.Secret{Name: "test", Data: []byte("secret")}))
			assert.NoError(t, s.Remove(model.Secret{Name: "test"}))
			assert.NoError(t, s.Update(model.Secret{Name: "bla", Type: model.PasswordType, Data: []byte("secret")}))
			assert.NotNil(t, s.Get(model.Secret{Name: "bla"}))
			assert.Nil(t, s.Get(model.Secret{Name: "test"}))
			assert.NotEmpty(t, s.ListKeys())
			assert.NotEmpty(t, s.ListKeys(model.PasswordType))

			assert.NoError(t, s.Close())
			assert.Nil(t, s.ListKeys())
			assert.Error(t, s.Remove(model.Secret{Name: "test"}))
			assert.Error(t, s.Update(model.Secret{Name: "bla", Data: []byte("secret")}))
		})
	}
}
