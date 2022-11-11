package client

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/gopherlearning/gophkeeper/internal"
	"github.com/rs/zerolog/log"
)

type CliCmd struct {
	db Repository
}

type CliState struct {
}

func (s *CliCmd) Run(ctx *internal.Context) error {
	log.Debug().
		Msg("я консольная часть менеджера паролей GophKeeper")
	SaveTermState()

	defer RestoreTermState()

	p := prompt.New(
		s.executor,
		s.completer,
		prompt.OptionShowCompletionAtStart(),
		prompt.OptionCompletionOnDown(),
		prompt.OptionPrefix("> "),
	)
	p.Run()

	return nil
}

var (
	rShow     = regexp.MustCompile(`^show ([\w/\.@-]+)$`)
	rEdit     = regexp.MustCompile(`^edit ([\w/\.@-]+)$`)
	rText     = regexp.MustCompile(`^text ([\w/\.@-]+)$`)
	rGenerate = regexp.MustCompile(`^generate ([\w/\.@-]+)$`)
	rBinary   = regexp.MustCompile(`^binary ([\w/\.@-]+)$`)
	rRemove   = regexp.MustCompile(`^remove ([\w/\.@-]+)$`)
)

func (s *CliCmd) executor(in string) {
	switch {
	case rShow.FindStringSubmatch(in) != nil:
		match := rShow.FindStringSubmatch(in)[1]
		fmt.Println(s.showText(match))
	case rText.FindStringSubmatch(in) != nil:
	case rGenerate.FindStringSubmatch(in) != nil:
	case rBinary.FindStringSubmatch(in) != nil:
	case rRemove.FindStringSubmatch(in) != nil:

	}
}

func (s *CliCmd) completer(in prompt.Document) []prompt.Suggest {
	switch {
	case !strings.Contains(in.CurrentLineBeforeCursor(), " "):
		return []prompt.Suggest{
			{Text: "show", Description: "Отобразить секрет"},
			{Text: "edit", Description: "Редактировать секрет"},
			{Text: "text", Description: "Создать текстовый секрет"},
			{Text: "generate", Description: "Сгенерировать пароль для текстового секрета"},
			{Text: "binary", Description: "Создать бинарный секрет"},
			{Text: "remove", Description: "Удалить секрет"},
		}
	case strings.Contains(in.TextBeforeCursor(), "show "):
		return s.keys()
	case strings.Contains(in.TextBeforeCursor(), "text "):
		return []prompt.Suggest{}
	}

	// fmt.Println(in.GetWordBeforeCursor())
	return []prompt.Suggest{
		{Text: "y", Description: "Применить введённую фразу"},
		{Text: "n", Description: "Ввести снова"},
	}
}

func (s *CliCmd) keys() []prompt.Suggest {
	res := make([]prompt.Suggest, 0)
	keys := s.db.ListKeys()

	for k := range keys {
		res = append(res, prompt.Suggest{Text: keys[k]})
	}

	return res
}

func (s *CliCmd) showText(key string) string {
	
}
