package client

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	errTest = fmt.Errorf("test error")
)

func Test_checkStorage(t *testing.T) {
	succsessPath := t.TempDir()
	require.NoError(t, os.Mkdir(succsessPath+"/.gophkeeper", 0755))
	failPath := t.TempDir()
	require.NoError(t, os.WriteFile(failPath+"/.gophkeeper", []byte{}, 0600))

	tests := []struct {
		err      error
		resError error
		name     string
		path     string
		wantErr  bool
	}{
		{
			name:     "first error",
			err:      errTest,
			wantErr:  true,
			resError: errTest,
		},
		{
			name:     "wrong storage path",
			err:      nil,
			wantErr:  true,
			path:     "/tmp/aaaaa",
			resError: ErrLocalStorageNotFound,
		},

		{
			name:     "not dir",
			err:      nil,
			wantErr:  true,
			path:     failPath,
			resError: ErrLocalStorageWrongType,
		},
		{
			name:     "success",
			err:      nil,
			wantErr:  false,
			path:     succsessPath,
			resError: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkStorage(tt.path, tt.err)
			if tt.wantErr {
				assert.ErrorIs(t, err, tt.resError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
