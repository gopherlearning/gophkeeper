package client

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gopherlearning/gophkeeper/internal"
	"github.com/gopherlearning/gophkeeper/internal/storage/local"
	"github.com/rs/zerolog/log"
	"golang.org/x/term"
)

// var termState *term.State

type Cmd struct {
	termState *term.State
	// Bla string
	// Authorize AuthorizeCmd `cmd:"" help:"Авторизоваться с помощью адреса сервера и мнемонической фразы и задать пароль для защиты локальной версии хранилища"`
}

// Run клиентская часть менеджера паролей GophKeeper.
func (l *Cmd) Run(ctx *internal.Context) error {
	log.Debug().
		Msg("я клиентская часть менеджера паролей GophKeeper")
	log.Debug().
		Msg("проверка хранилища")

	var serverURL string

	err := func(err error) error {
		if err != nil {
			if errors.Is(err, ErrLocalStorageNotFound) {
				ac := &AuthorizeCmd{cmd: l}
				if err = ac.Run(ctx); err != nil {
					return err
				}

				serverURL = ac.SeverURL

				return nil
			}

			return err
		}

		return nil
	}(checkStorage(os.UserConfigDir()))

	if err != nil {
		return err
	}

	path, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	key, err := os.ReadFile(filepath.Join(path, ".gophkeeper", "CACHE"))
	if err != nil {
		return err
	}

	db, err := local.NewLocalStorage(fmt.Sprint(key)[10:42], filepath.Join(path, ".gophkeeper"), serverURL)
	if err != nil {
		return err
	}

	cli := &CliCmd{db: db, cmd: l}

	return cli.Run(ctx)
}

// SaveTermState - save terminal state on start.
func (l *Cmd) SaveTermState() {
	oldState, err := term.GetState(int(os.Stdin.Fd()))
	if err != nil {
		return
	}

	l.termState = oldState
}

// RestoreTermState - restore terminal state on exit.
func (l *Cmd) RestoreTermState() {
	if r := recover(); r != nil {
		log.Error().Msgf("Recovered in f %v", r)
	}

	if l.termState != nil {
		err := term.Restore(int(os.Stdin.Fd()), l.termState)
		if err != nil {
			log.Error().Msgf("Recovered in f.Error: %v", err)
		}
	}
}
