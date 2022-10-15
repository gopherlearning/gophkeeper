package conf

import (
	"github.com/gopherlearning/gophkeeper/internal/client"
	"github.com/gopherlearning/gophkeeper/internal/server"
)

type Args struct {
	Version   VersionCmd          `cmd:"" help:"Показать информацию о версии"`
	Server    server.Cmd          `cmd:"" help:"API сервер для менеджера паролей GophKeeper"`
	Authorize client.AuthorizeCmd `cmd:"" help:"Авторизоваться с помощью адреса сервера и мнемонической фразы и задать пароль для защиты локальной версии хранилища"`
	Verbose   bool                `name:"verbose" short:"v" help:"Включить расширенное логирование"`
	// Client    client.Cmd          `cmd:"" help:"клиентская часть менеджера паролей GophKeeper" default:"withargs"`
}

func (a *Args) Debug() bool {
	return a.Verbose
}
