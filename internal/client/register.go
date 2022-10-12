package client

import (
	"github.com/gopherlearning/gophkeeper/internal"
	"github.com/rs/zerolog/log"
)

type RegisterCmd struct {
	SeverURL string `name:"url" help:"Адрес API сервера"`
}

func (l *RegisterCmd) Run(ctx *internal.Context) error {
	log.Info().
		Msg("тут будет запущен процесс регистрации")

	return nil
}
