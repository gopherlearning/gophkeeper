package conf

import "github.com/rs/zerolog/log"

type VersionCmd struct {
}

func (l *VersionCmd) Run(ctx *Context) error {
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
