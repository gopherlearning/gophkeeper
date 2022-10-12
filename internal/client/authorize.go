package client

import (
	"fmt"

	"github.com/gopherlearning/gophkeeper/internal"
	"github.com/rs/zerolog/log"
)

type AuthorizeCmd struct {
	SeverURL string `name:"url" help:"Адрес API сервера"`
}

func (l *AuthorizeCmd) Run(ctx *internal.Context) error {
	log.Info().
		Msg("тут будет запущен процесс авторизации")

	fmt.Println(2)

	return nil
}
