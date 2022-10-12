package client

import (
	"fmt"

	"github.com/gopherlearning/gophkeeper/internal"
	"github.com/rs/zerolog/log"
)

type Cmd struct {
	Register  RegisterCmd  `cmd:"" help:"Указать адрес API сервера, задать пароль и сгенерировать мнемоническую фразу (пароль защищает локальную версию хранилища)"`
	Authorize AuthorizeCmd `cmd:"" help:"Авторизоваться с помощью адреса сервера и мнемонической фразы и задать пароль для защиты локальной версии хранилища"`
}

// Run клиентская часть менеджера паролей GophKeeper.
func (l *Cmd) Run(ctx *internal.Context) error {
	fmt.Println(1)

	log.Info().
		Msg("я клтентская часть менеджера паролей GophKeeper")

	return nil
}
