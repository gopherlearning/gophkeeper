package client

import (
	"fmt"
	"regexp"

	"github.com/c-bata/go-prompt"
	"github.com/eiannone/keyboard"
	"github.com/gopherlearning/gophkeeper/internal"
	"github.com/rs/zerolog/log"
)

var serverURLRegexp = regexp.MustCompile(`^grpc[s]?://[\w/\.@-]+:[0-9]+$`)

type ServerURL struct {
	cmd *Cmd
	URL string
}

func (l *ServerURL) Run(ctx *internal.Context) error {
	fmt.Print("Планируете ли использовать удалённый сервер? Да(y) / Нет(n)")

	for {
		char, _, err := keyboard.GetSingleKey()
		if err != nil {
			return err
		}

		switch char {
		case 'y', 'Y', 'н', 'Н':
			l.cmd.SaveTermState()
			defer l.cmd.RestoreTermState()
			fmt.Println(" - y\nВведите URL сервера (grpc[s]://<address>:<port>):")

			p := prompt.New(
				l.Executor,
				l.Completer,
				prompt.OptionCompletionOnDown(),
				prompt.OptionPrefix("> "),
				prompt.OptionSetExitCheckerOnInput(func(in string, breakline bool) bool {
					return in == "y"
				}),
			)
			p.Run()
			log.Debug().Msg("ввод адреса закончен")

			return nil
		case 'n', 'N', 'т', 'Т':
			return nil
		}
	}
}

func (l *ServerURL) Executor(in string) {
	if in == "n" {
		l.URL = ""
		return
	}

	if serverURLRegexp.FindStringSubmatch(in) == nil {
		fmt.Println("Введите корректный адрес сервера")
		return
	}

	l.URL = serverURLRegexp.FindStringSubmatch(in)[0]
}

func (l *ServerURL) Completer(in prompt.Document) []prompt.Suggest {
	if len(l.URL) == 0 {
		return []prompt.Suggest{}
	}

	return []prompt.Suggest{
		{Text: "y", Description: "Применить"},
		{Text: "n", Description: "Не применять"},
	}
}
