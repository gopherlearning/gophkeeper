package server

import (
	"github.com/gopherlearning/gophkeeper/internal"
	"github.com/rs/zerolog/log"
)

type Cmd struct {
}

func (l *Cmd) Run(ctx *internal.Context) error {
	log.Info().
		Msg("я API сервер для менеджера паролей GophKeeper")
	return nil
}
