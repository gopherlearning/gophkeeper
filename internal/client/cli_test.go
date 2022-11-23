package client

import (
	"testing"

	"github.com/c-bata/go-prompt"
	"github.com/stretchr/testify/assert"
)

func Test_pathDiscover(t *testing.T) {
	tests := []struct {
		name string
		dir  string
		want []prompt.Suggest
	}{
		{
			name: "success",
			dir:  t.TempDir(),
			want: []prompt.Suggest{{Text: "001"}},
		},
		{
			name: "error",
			dir:  "tmp////ddddd",
			want: []prompt.Suggest{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := pathDiscover(tt.dir)
			assert.ElementsMatch(t, tt.want, got)
		})
	}
}
