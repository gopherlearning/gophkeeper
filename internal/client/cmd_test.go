package client

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSaveTermState(t *testing.T) {
	t.Run("not pannic", func(t *testing.T) {
		cmd := &Cmd{}
		require.NotPanics(t, func() { cmd.SaveTermState() })
		fmt.Println(cmd.termState)
		require.NotPanics(t, func() { cmd.RestoreTermState() })
	})
}
