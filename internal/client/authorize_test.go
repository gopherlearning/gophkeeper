package client

import (
	"os"
	"testing"

	"github.com/c-bata/go-prompt"
	"github.com/gopherlearning/gophkeeper/internal"
	"github.com/stretchr/testify/assert"
)

func Test_generateMnemonic(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, err := generateMnemonic()
		assert.NoError(t, err)
		assert.NotEmpty(t, m)
	})
	t.Run("NewEntropy error", func(t *testing.T) {
		entropyBitSize = 127
		m, err := generateMnemonic()
		assert.Error(t, err)
		assert.Empty(t, m)
		entropyBitSize = 128
	})
	t.Run("NewMnemonic error", func(t *testing.T) {
		internal.EmulatedError = "NewMnemonic error 1"
		m, err := generateMnemonic()
		assert.Error(t, err)
		assert.Empty(t, m)
		internal.EmulatedError = ""
	})
	t.Run("compiler", func(t *testing.T) {
		for _, v := range authorizeTests {
			if len(v.res) != 0 {
				assert.NotContains(t, v.state.Words(), v.state.Completer(v.pr))
			}
			assert.NotEmpty(t, v.state.Words())
		}
	})
	t.Run("executor", func(t *testing.T) {
		os.Setenv("HOME", t.TempDir())
		for _, v := range authorizeTests {
			if v.pr.Text != "y" {
				v.state.Executor(v.pr.Text)
			}
			if v.mustReady {
				assert.True(t, v.state.ready)
			}
			if v.pr.Text == "y" {
				assert.Panics(t, func() { v.state.Executor(v.pr.Text) })
			}
			assert.NotEmpty(t, v.state.Words())
		}
	})
}

var authorizeTests = []struct {
	cmd       *AuthorizeCmd
	state     *AuthorizeState
	res       []prompt.Suggest
	pr        prompt.Document
	mustReady bool
}{
	{
		state: &AuthorizeState{},
		res:   []prompt.Suggest{},
		pr:    prompt.Document{},
	},
	{
		state: &AuthorizeState{words: "abstract absurd abuse access accident account accuse achieve acid acoustic acquire across", ready: true},
		res: []prompt.Suggest{
			{
				Text:        "y",
				Description: "Применить введённую фразу",
			},
			{
				Text:        "n",
				Description: "Ввести снова",
			},
		},
		pr: prompt.Document{Text: "n"},
	},
	{
		state: &AuthorizeState{words: "abstract absurd abuse access accident account accuse achieve acid acoustic acquire across", ready: true},
		res: []prompt.Suggest{
			{
				Text:        "y",
				Description: "Применить введённую фразу",
			},
			{
				Text:        "n",
				Description: "Ввести снова",
			},
		},
		pr: prompt.Document{Text: "y"},
	},
	{
		state:     &AuthorizeState{words: "abstract absurd abuse access accident account accuse achieve acid acoustic acquire"},
		res:       []prompt.Suggest{{Text: "abstract"}, {Text: "absurd"}, {Text: "abuse"}, {Text: "access"}, {Text: "accident"}, {Text: "account"}, {Text: "accuse"}, {Text: "achieve"}, {Text: "acid"}, {Text: "acoustic"}, {Text: "acquire"}, {Text: "across"}},
		pr:        prompt.Document{Text: "abstract absurd abuse access accident account accuse achieve acid acoustic acquire across "},
		mustReady: true,
	},
	{
		state: &AuthorizeState{words: ""},
		res:   []prompt.Suggest{{Text: "abstract"}, {Text: "absurd"}, {Text: "abuse"}, {Text: "access"}, {Text: "accident"}, {Text: "account"}, {Text: "accuse"}, {Text: "achieve"}, {Text: "acid"}, {Text: "acoustic"}, {Text: "acquire"}},
		pr:    prompt.Document{Text: "abstract absurd abuse access accident account accuse achieve acid acoustic acquire"},
	},
	{
		state: &AuthorizeState{words: "abstract blabla absurd abuse access accident account accuse achieve acid acoustic acquire"},
		res:   []prompt.Suggest{{Text: "abstract"}, {Text: "absurd"}, {Text: "abuse"}, {Text: "access"}, {Text: "accident"}, {Text: "account"}, {Text: "accuse"}, {Text: "achieve"}, {Text: "acid"}, {Text: "acoustic"}, {Text: "acquire"}},
		// pr:    prompt.Document{Text: "blabla absurd abuse access accident account accuse achieve acid acoustic acquire"},
	},
}
