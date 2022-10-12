package conf

import (
	"github.com/gopherlearning/gophkeeper/internal"
	"github.com/rs/zerolog/log"
)

type VersionCmd struct {
}

func (l *VersionCmd) Run(ctx *internal.Context) error {
	log.Info().
		Str("buildVersion", buildVersion).
		Str("buildDate", buildDate).
		Str("buildCommit", buildCommit).
		Msg("")

	return nil
}

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)
