package client

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/gopasspw/gopass/pkg/pwgen"
	"github.com/gopherlearning/gophkeeper/internal"
	"github.com/gopherlearning/gophkeeper/internal/model"
	"github.com/gopherlearning/gophkeeper/pkg/editor"
	"github.com/rs/zerolog/log"
)

type CliCmd struct {
	temp  *model.Secret
	cmd   *Cmd
	db    Repository
	exist bool
}

type CliState struct {
}

var (
	baseCommands = []prompt.Suggest{
		{Text: "show", Description: "Отобразить секрет"},
		{Text: "edit", Description: "Редактировать секрет"},
		{Text: "text", Description: "Создать текстовый секрет"},
		{Text: "card", Description: "Создать банковскую карту"},
		{Text: "password", Description: "Создать пару логин/пароль"},
		{Text: "generate", Description: "Сгенерировать пароль (not implemented)"},
		{Text: "binary", Description: "Создать бинарный секрет"},
		{Text: "remove", Description: "Удалить секрет"},
	}
)

func (s *CliCmd) Run(ctx *internal.Context) error {
	log.Debug().
		Msg("хранилиже подгружено")
	s.cmd.SaveTermState()

	defer s.cmd.RestoreTermState()

	p := prompt.New(
		s.executor,
		s.completer,
		prompt.OptionShowCompletionAtStart(),
		prompt.OptionCompletionOnDown(),
		prompt.OptionLivePrefix(s.db.Status()),
		prompt.OptionPrefix("❌ >"),
	)
	p.Run()

	return nil
}

var (
	rShow     = regexp.MustCompile(`^show ([\w/\.@-]+)$`)
	rEdit     = regexp.MustCompile(`^edit ([\w/\.@-]+)$`)
	rText     = regexp.MustCompile(`^text ([\w/\.@-]+)$`)
	rCard     = regexp.MustCompile(`^card ([\w/\.@-]+)$`)
	rPassword = regexp.MustCompile(`^password ([\w/\.@-]+)$`)
	rGenerate = regexp.MustCompile(`^generate ([\w/\.@-]+)$`)
	rBinary   = regexp.MustCompile(`^binary ([\w/\.@-]+)( )?([\w/\.-@_]+)?$`)
	rRemove   = regexp.MustCompile(`^remove ([\w/\.@-]+)$`)
)

func (s *CliCmd) executor(in string) {
	switch {
	case rShow.FindStringSubmatch(in) != nil:
		match := rShow.FindStringSubmatch(in)[1]
		fmt.Println(s.showText(match))
	case rText.FindStringSubmatch(in) != nil:
		s.text(rText.FindStringSubmatch(in)[1])
	case rCard.FindStringSubmatch(in) != nil:
		s.card(rCard.FindStringSubmatch(in)[1])
	case rPassword.FindStringSubmatch(in) != nil:
	case rEdit.FindStringSubmatch(in) != nil:
		s.edit(rEdit.FindStringSubmatch(in)[1])
	case rGenerate.FindStringSubmatch(in) != nil:
		s.generate(rGenerate.FindStringSubmatch(in)[1])
	case s.exist:
		if in == "y" {
			s.generate("")
		}

		s.exist = false
		s.temp = nil
	case rBinary.FindStringSubmatch(in) != nil:
		s.binary(rBinary.FindStringSubmatch(in)[1:])
	case rRemove.FindStringSubmatch(in) != nil:
		err := s.db.Remove(model.Secret{Name: rRemove.FindStringSubmatch(in)[1]})
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (s *CliCmd) completer(in prompt.Document) []prompt.Suggest {
	switch {
	case !strings.Contains(in.CurrentLineBeforeCursor(), " ") && s.exist:
		return []prompt.Suggest{
			{Text: "y", Description: "заменить секрет"},
			{Text: "n", Description: "не заменять секрет"},
		}
	case !strings.Contains(in.CurrentLineBeforeCursor(), " ") && !s.exist:
		return prompt.FilterHasPrefix(baseCommands, in.GetWordBeforeCursor(), false)
	case
		strings.Contains(in.TextBeforeCursor(), "show "),
		strings.Contains(in.TextBeforeCursor(), "edit "),
		strings.Contains(in.TextBeforeCursor(), "remove "):
		return prompt.FilterHasPrefix(s.keys(), in.GetWordBeforeCursor(), false)
	case strings.Contains(in.TextBeforeCursor(), "generate "):
		return prompt.FilterHasPrefix(s.keys(model.PasswordType), in.GetWordBeforeCursor(), false)
	case strings.Contains(in.TextBeforeCursor(), "text "):
		return prompt.FilterHasPrefix(s.keys(model.TextType), in.GetWordBeforeCursor(), false)
	case strings.Contains(in.TextBeforeCursor(), "card "):
		return prompt.FilterHasPrefix(s.keys(model.CardType), in.GetWordBeforeCursor(), false)
	case strings.Contains(in.TextBeforeCursor(), "password "):
		return prompt.FilterHasPrefix(s.keys(model.PasswordType), in.GetWordBeforeCursor(), false)
	case strings.Contains(in.TextBeforeCursor(), "binary "):
		return s.binaryCase(in)
	default:
		return []prompt.Suggest{}
	}
}

func (s *CliCmd) binaryCase(in prompt.Document) []prompt.Suggest {
	str := rBinary.FindStringSubmatch(in.TextBeforeCursor())
	if len(str) == 0 {
		return prompt.FilterHasPrefix(s.keys(model.BinaryType), in.GetWordBeforeCursor(), false)
	}

	if len(str[2]) != 0 {
		return prompt.FilterHasPrefix(pathDiscover(str[3]), in.GetWordBeforeCursor(), false)
	}

	return prompt.FilterHasPrefix(s.keys(model.BinaryType), in.GetWordBeforeCursor(), false)
}

func (s *CliCmd) keys(types ...model.SecretType) []prompt.Suggest {
	res := make([]prompt.Suggest, 0)
	keys := s.db.ListKeys(types...)

	for k := range keys {
		res = append(res, prompt.Suggest{Text: keys[k]})
	}

	return res
}

func (s *CliCmd) showText(key string) string {
	secret := s.db.Get(model.Secret{Name: key})
	if secret == nil {
		return "нет такого секрета"
	}

	return secret.String()
}

func (s *CliCmd) generate(name string) {
	if s.temp != nil {
		err := s.db.Update(*s.temp)
		if err != nil {
			fmt.Println(err)
		}

		s.temp = nil

		return
	}

	secret := &model.Secret{
		Name: name,
		Type: model.PasswordType,
		Data: []byte(pwgen.GeneratePassword(24, false)),
	}

	if s.db.Get(model.Secret{Name: name}) != nil {
		s.exist = true
		s.temp = secret

		fmt.Println("секрет существует, заменить?")

		return
	}

	err := s.db.Update(*secret)
	if err != nil {
		fmt.Println(err)
	}

	s.temp = nil
}

func (s *CliCmd) edit(name string) {
	secret := s.db.Get(model.Secret{Name: name})
	if secret == nil {
		secret = &model.Secret{
			Name: name,
			Data: []byte{},
			Type: model.TextType,
		}
	}

	edited, err := editor.Edit(secret.Data)
	if err != nil {
		fmt.Printf("ошибка редактирования секрета: %v\n", err)
		return
	}

	secret.Data = edited
	err = s.db.Update(*secret)

	if err != nil {
		fmt.Printf("ошибка обновления секрета: %v\n", err)
		return
	}
}

func (s *CliCmd) card(name string) {
	secret := s.db.Get(model.Secret{Name: name})
	if secret != nil {
		if secret.Type != model.CardType {
			fmt.Println("секрет не является данными банковской карты")
			return
		}

		s.edit(name)

		return
	}

	edited, err := editor.Edit([]byte(model.CadrTemplate))
	if err != nil {
		fmt.Printf("ошибка редактирования секрета: %v\n", err)
		return
	}

	if len(edited) == 0 || len(edited) == len([]byte(model.CadrTemplate)) {
		return
	}

	err = s.db.Update(model.Secret{
		Name: name,
		Data: edited,
		Type: model.CardType,
	})
	if err != nil {
		fmt.Printf("ошибка обновления секрета: %v\n", err)
	}
}

func (s *CliCmd) text(name string) {
	secret := s.db.Get(model.Secret{Name: name})
	if secret != nil {
		if secret.Type != model.TextType {
			fmt.Println("секрет не является текстовыми данными")
			return
		}

		s.edit(name)

		return
	}

	edited, err := editor.Edit([]byte{})
	if err != nil {
		fmt.Printf("ошибка редактирования секрета: %v\n", err)
		return
	}

	if len(edited) == 0 {
		return
	}

	err = s.db.Update(model.Secret{
		Name: name,
		Data: edited,
		Type: model.TextType,
	})
	if err != nil {
		fmt.Printf("ошибка обновления секрета: %v\n", err)
	}
}

func (s *CliCmd) binary(nameAndPath []string) {
	if len(nameAndPath) < 1 {
		return
	}

	secret := s.db.Get(model.Secret{Name: nameAndPath[0]})
	if secret != nil {
		if secret.Type != model.BinaryType {
			fmt.Println("секрет не является бинарным")
			return
		}

		if len(nameAndPath) != 3 || len(nameAndPath[2]) == 0 {
			fmt.Println("не указан путь для сохранения")
			return
		}

		err := os.WriteFile(nameAndPath[2], secret.Data, 0600)
		if err != nil {
			fmt.Printf("не удалось извлечь бинарные данные: %v\n", err)
			return
		}

		return
	}

	if len(nameAndPath) != 3 || len(nameAndPath[2]) == 0 {
		fmt.Println("не указан путь до файла")
		return
	}

	file, err := os.ReadFile(nameAndPath[2])
	if err != nil {
		fmt.Printf("не удалось добавить бинарные данные: %v\n", err)
		return
	}

	err = s.db.Update(model.Secret{Name: nameAndPath[0], Type: model.BinaryType, Data: file})
	if err != nil {
		fmt.Printf("не удалось добавить бинарные данные: %v\n", err)
	}
}

func pathDiscover(dir string) []prompt.Suggest {
	res := make([]prompt.Suggest, 0)
	dirs, err := os.ReadDir(path.Dir(dir))

	if err != nil {
		return []prompt.Suggest{}
	}

	for _, v := range dirs {
		res = append(res, prompt.Suggest{Text: v.Name()})
	}

	return res
}
