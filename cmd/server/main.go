package main

import (
	"github.com/gopherlearning/gophkeeper/internal/conf"
	"github.com/rs/zerolog"
)

var (
	cfg conf.Args
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
}

func main() {
	conf.Run("gophkeeper-server", "API сервер для менеджера паролей GophKeeper", &cfg)
}
