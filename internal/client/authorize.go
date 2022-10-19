package client

import (
	"fmt"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/cosmos/go-bip39"
	"github.com/eiannone/keyboard"
	"github.com/gopherlearning/gophkeeper/internal"
)

var (
	entropyBitSize = 128
)

type AuthorizeCmd struct {
	SeverURL string `name:"url" help:"Адрес API сервера"`
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
			SaveTermState()

			defer RestoreTermState()
			fmt.Println(" - y\nВведите мнемоническую фразу (12 слов):")

			s := &AuthorizeState{}
			p := prompt.New(
				s.executor,
				s.completer,
				prompt.OptionShowCompletionAtStart(),
				prompt.OptionCompletionOnDown(),
				prompt.OptionPrefix("> "),
			)
			p.Run()

			return nil
		case 'n', 'N', 'т', 'Т':
			fmt.Println(" - n\nЗапустить генерацию")
			return nil
		case 'e', 'E', 'q', 'Q', 'й', 'Й', 'у', 'У':
			fmt.Println()
			return nil
		}
	}
}

func (s *AuthorizeState) executor(in string) {
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

		fmt.Println(in, "\nВсё правильно?")

		s.ready = true
	}

	if strings.TrimSpace(in) == "n" {
		s.ready = false
		s.words = ""

		return
	}

	if strings.TrimSpace(in) == "y" {
		s.confirm = true
	}

	return
	// if s.ready {
	// 	fmt.Println(in, "\nВсё правильно?")
	// }

	// if s.confirm {
	// 	if len(s.password) < 8 {
	// 		fmt.Println(in, "\nВведите пароль для локального хранилища (минимум 8 символов):")
	// 	}
	// 	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	defer term.Restore(int(os.Stdin.Fd()), oldState)
	// 	s.password = in
	// }
}

func (s *AuthorizeState) completer(in prompt.Document) []prompt.Suggest {
	if s.confirm {
		return []prompt.Suggest{}
	}

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

type AuthorizeState struct {
	words    string
	password string
	ready    bool
	confirm  bool
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
