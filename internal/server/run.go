package server

import (
	"fmt"

	"github.com/gopherlearning/gophkeeper/internal"
	"github.com/rs/zerolog/log"
)

type Cmd struct {
	Config  string `short:"c" type:"existingfile" help:"Конфигурационный файл для сервера." default:"config.yaml"`
	Listen  string `yaml:"listen" help:"Порт GRPC сервера" default:":9765"`
	Storage struct {
		Path string `yaml:"path" help:"Путь к локальному хранилищу" default:"gopher-storage"`
		DSN  string `yaml:"dsn" help:"Адрес базы данных (не реализовано)"`
	} `embed:"" prefix:"storage." help:"Параметры хранилища"`
}

func (l *Cmd) Run(ctx *internal.Context) error {
	fmt.Printf("%+v\n", l)
	log.Info().
		Msg("я API сервер для менеджера паролей GophKeeper")

	return nil
}
