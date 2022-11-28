package client

import (
	"fmt"
	"os"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/cosmos/go-bip39"
	"github.com/eiannone/keyboard"
	"github.com/gopherlearning/gophkeeper/internal"
	"github.com/rs/zerolog/log"
)

var (
	entropyBitSize = 128
)

type AuthorizeCmd struct {
	cmd      *Cmd
	SeverURL string `name:"url" help:"Адрес API сервера"`
}
type AuthorizeState struct {
	cmd   *Cmd
	words string
	ready bool
}

func (l *AuthorizeCmd) Run(ctx *internal.Context) error {
	fmt.Print("Есть ли у Вас мнемоническая фраза? Да(y) / Нет(n) / Выход(q)")

	for {
		char, _, err := keyboard.GetSingleKey()
		if err != nil {
			return err
		}

		switch char {
		case 'y', 'Y', 'н', 'Н':
			l.cmd.SaveTermState()
			defer l.cmd.RestoreTermState()
			fmt.Println(" - y\nВведите мнемоническую фразу (12 слов):")

			s := &AuthorizeState{cmd: l.cmd}
			p := prompt.New(
				s.Executor,
				s.Completer,
				prompt.OptionShowCompletionAtStart(),
				prompt.OptionCompletionOnDown(),
				prompt.OptionSetExitCheckerOnInput(func(in string, breakline bool) bool { return in == "y" }),
				prompt.OptionPrefix("> "),
			)
			p.Run()
			log.Debug().Msg("ввод мнемонической фразы закончен")

			return l.initStorage(ctx, s.words)
		case 'n', 'N', 'т', 'Т':
			mnemonic, err := generateMnemonic()
			if err != nil {
				return err
			}

			fmt.Printf("\nСохраните данную фразу в надёжном месте:\n%s\n", mnemonic)

			return l.initStorage(ctx, mnemonic)
		case 'e', 'E', 'q', 'Q', 'й', 'Й', 'у', 'У':
			fmt.Println()
			return nil
		}
	}
}

func (l *AuthorizeCmd) initStorage(ctx *internal.Context, mnemonic string) error {
	serv := ServerURL{cmd: l.cmd}
	err := serv.Run(ctx)

	if err != nil {
		return err
	}

	l.SeverURL = serv.URL
	path, err := os.UserConfigDir()

	if err != nil {
		return err
	}

	return initStorage(mnemonic, path, serv.URL)
}

func (s *AuthorizeState) Executor(in string) {
	words := strings.Split(strings.TrimSpace(in), " ")
	if len(words) == 12 {
		for _, w := range words {
			finded := false

			for _, v := range bip39.EnglishWordList {
				if w == v {
					finded = true
					break
				}
			}

			if !finded {
				fmt.Printf("\nНедопустимое слово: %s. Введите снова\n", w)

				s.ready = false
				s.words = ""

				return
			}
		}

		fmt.Println("Всё верно??")

		s.ready = true
	}

	if strings.TrimSpace(in) == "n" {
		s.ready = false
		s.words = ""

		return
	}

	if strings.TrimSpace(in) == "y" {
		return
		// path, err := os.UserConfigDir()
		// if err != nil {
		// 	s.cmd.RestoreTermState()
		// 	log.Err(err)
		// 	os.Exit(1)
		// }

		// var serverURL string

		// serv := &ServerURL{}

		// serv.Run(&internal.Context{})

		// err = initStorage(s.words, path, serverURL)
		// if err != nil {
		// 	log.Err(err)
		// }

		// s.cmd.RestoreTermState()
		// os.Exit(0)
	}
}

func (s *AuthorizeState) Completer(in prompt.Document) []prompt.Suggest {
	if s.ready {
		return []prompt.Suggest{
			{Text: "y", Description: "Применить введённую фразу"},
			{Text: "n", Description: "Ввести снова"},
		}
	}

	ins := make([]string, 0)

	for _, v := range strings.Split(in.CurrentLine(), " ") {
		if len(v) != 0 {
			ins = append(ins, v)
		}
	}

	s.words = strings.Join(ins, " ")

	if len(ins) >= 12 {
		return []prompt.Suggest{}
	}

	return prompt.FilterHasPrefix(s.Words(), in.GetWordBeforeCursor(), false)
}

// generateMnemonic сгенерировать мнемоническую фразу.
func generateMnemonic() (mnemonic string, err error) {
	entropy, err := bip39.NewEntropy(entropyBitSize)
	if err != nil {
		return "", err
	}

	mnemonic, err = bip39.NewMnemonic(entropy)

	if err = internal.EmulateError(err, 1); err != nil {
		return "", err
	}

	return mnemonic, nil
}

func (s *AuthorizeState) Words() []prompt.Suggest {
	m := make([]prompt.Suggest, 0)

	for _, v := range bip39.EnglishWordList {
		if strings.Contains(s.words, v) {
			continue
		}

		m = append(m, prompt.Suggest{Text: v})
	}

	return m
}
