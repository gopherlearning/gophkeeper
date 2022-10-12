package main

import (
	"os"

	"github.com/gopherlearning/gophkeeper/internal/conf"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	cfg conf.Args
)

func init() {
	log.Logger = zerolog.
		New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05"}).
		With().Timestamp().Logger()
}

func main() {
	conf.Run("gophkeeper", "Менеджер паролей GophKeeper", &cfg)
}
