package conf

import (
	"os"
	"strconv"
	"strings"

	"github.com/alecthomas/kong"
	kongyaml "github.com/alecthomas/kong-yaml"
	"github.com/gopherlearning/gophkeeper/internal"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Run задаёт стандартные значения, читает конфиг, проверяет переменное окружение и флаги.
func Run(name, desc string, cfg cmdContext) {
	var confPath string

	for i := 0; i < len(os.Args); i++ {
		if os.Args[i] == "-c" || os.Args[i] == "--config" {
			confPath = os.Args[i+1]
		}
	}

	opts := []kong.Option{
		kong.Name(name),
		kong.Description(desc),
		kong.UsageOnError(),
	}

	if len(confPath) != 0 {
		opts = append(opts, kong.Configuration(kongyaml.Loader, confPath))
	}

	parser := kong.Must(cfg, opts...)

	ctx, err := parser.Parse(os.Args[1:])
	parser.FatalIfErrorf(err)
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if cfg.Debug() {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)

		log.Logger = log.With().
			Caller().
			Logger()

		zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
			return strings.Split(file, "gophkeeper/")[1] + ":" + strconv.Itoa(line)
		}
	}

	err = ctx.Run(&internal.Context{Verbose: cfg.Debug()})

	ctx.FatalIfErrorf(err)
}

type cmdContext interface {
	Debug() bool
}
