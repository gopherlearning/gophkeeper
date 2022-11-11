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

var termState *term.State

type Cmd struct {
	// Bla string
	// Authorize AuthorizeCmd `cmd:"" help:"Авторизоваться с помощью адреса сервера и мнемонической фразы и задать пароль для защиты локальной версии хранилища"`
}

// Run клиентская часть менеджера паролей GophKeeper.
func (l *Cmd) Run(ctx *internal.Context) error {
	log.Debug().
		Msg("я клтентская часть менеджера паролей GophKeeper")

	err := checkStorage(os.UserConfigDir())
	if err != nil {
		if errors.Is(err, ErrLocalStorageNotFound) {
			return new(AuthorizeCmd).Run(ctx)
		}

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

	db, err := local.NewLocalStorage(fmt.Sprint(key)[10:42], filepath.Join(path, ".gophkeeper"), nil)
	if err != nil {
		return err
	}

	cli := &CliCmd{db: db}

	return cli.Run(ctx)
}

// SaveTermState - save terminal state on start.
func SaveTermState() {
	oldState, err := term.GetState(int(os.Stdin.Fd()))
	if err != nil {
		return
	}

	termState = oldState
}

// RestoreTermState - restore terminal state on exit.
func RestoreTermState() {
	if r := recover(); r != nil {
		fmt.Println("Recovered in f", r)
	}

	if termState != nil {
		err := term.Restore(int(os.Stdin.Fd()), termState)
		if err != nil {
			fmt.Println("Recovered in f.Error: ", err)
		}
	}
}
