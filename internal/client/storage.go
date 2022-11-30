package client

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gopherlearning/gophkeeper/internal/storage/local"
	"github.com/rs/zerolog/log"
)

var (
	ErrLocalStorageNotFound  = errors.New("storage does not exist")
	ErrLocalStorageWrongType = errors.New("wrong storage type")
)

func checkStorage(path string, err error) error {
	if err != nil {
		return err
	}

	path = filepath.Join(path, ".gophkeeper")
	st, err := os.Stat(path)

	if err != nil {
		return ErrLocalStorageNotFound
	}

	if !st.IsDir() {
		return ErrLocalStorageWrongType
	}

	return nil
}

func initStorage(mnemonic, path, serverURL string) error {
	fmt.Println("Инициализация ")

	path = filepath.Join(path, ".gophkeeper")
	s := sha256.Sum256([]byte(fmt.Sprint(`%%%%`, mnemonic)))

	log.Debug().Msg(fmt.Sprint(s)[10:42])
	db, err := local.NewLocalStorage(fmt.Sprint(s)[10:42], path, serverURL)

	if err != nil {
		return err
	}
	defer db.Close()

	err = os.WriteFile(filepath.Join(path, "CACHE"), s[:], 0600)
	if err != nil {
		return err
	}

	fmt.Println("Хранилище менеджера паролей создано.")

	return nil
}
