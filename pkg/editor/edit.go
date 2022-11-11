package editor

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/rs/zerolog/log"
)

var editor = "nano"

func init() {
	if s := os.Getenv("EDITOR"); s != "" {
		editor = s
	}
}

func Edit(origin []byte) ([]byte, error) {
	tmp, err := os.CreateTemp(os.TempDir(), "pkg-editor")
	if err != nil {
		return nil, err
	}

	err = tmp.Close()
	if err != nil {
		return nil, err
	}

	defer func() {
		err = os.Remove(tmp.Name())
		if err != nil {
			log.Err(err)
		}
	}()

	err = os.WriteFile(tmp.Name(), origin, 0600)
	if err != nil {
		return nil, err
	}

	path, err := exec.LookPath(editor)
	if err != nil {
		return nil, fmt.Errorf("error %s while looking up for %s", path, editor)
	}

	cmd := exec.Command(path, tmp.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	err = cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("start failed: %v", err)
	}

	err = cmd.Wait()
	b, err := os.ReadFile(tmp.Name())

	if err != nil {
		return nil, fmt.Errorf("failed read file %s: %v", tmp.Name(), err)
	}

	return b, nil
}
