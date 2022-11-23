package conf

import (
	"github.com/gopherlearning/gophkeeper/internal/client"
	"github.com/gopherlearning/gophkeeper/internal/server"
)

type Args struct {
	Cli     client.Cmd `cmd:"" help:"клиентская часть менеджера паролей GophKeeper" default:"1"`
	Version VersionCmd `cmd:"" help:"Показать информацию о версии"`
	Server  server.Cmd `cmd:"" help:"API сервер для менеджера паролей GophKeeper"`
	Verbose bool       `name:"verbose" short:"v" help:"Включить расширенное логирование"`
}

func (a *Args) Debug() bool {
	return a.Verbose
}
