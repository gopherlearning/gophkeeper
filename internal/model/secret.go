package model

import (
	"bytes"
	"fmt"
	"text/template"
)

var secretTemplate = template.Must(template.New("secret").Parse(`---
{{.Text}}
---{{if .S.Labels}}
labels:{{range $index, $element := .S.Labels}}
  {{$index}}: {{$element}}{{end}}{{end}}
name: {{.S.Name}}
owner: {{.S.Owner}}
`))

// Secret является минимальной единицей хранения, должен содержать Owner, Name
// может содержать Labels.
type Secret struct {
	Name   string
	Owner  string
	Labels map[string]string
	data   []byte
	Type   SecretType
}

type aliasSecret struct {
	S *Secret
	// d Decryptor
	// e Encryptor
}

// String возвращает расшифрованное текстовое представление сожердимого секрета,
/*
---
password: dededededede
la:
  la: 1
plaintext
---
name: supersecret
labels:
  readers: [fefefefefe, tgtgtgtg]
  editors: [eeddddeee]
*/
func (s *Secret) String() string {
	out := &bytes.Buffer{}
	_ = secretTemplate.Execute(out, aliasSecret{S: s})

	return out.String()
}

// Text возвращает расшифрованное сожердимое переменной s.Data.
func (s *Secret) Text() string {
	if s.Type != TextType {
		return fmt.Errorf("это не текстовые данные. Тип - %s", s.Type).Error()
	}

	return string(s.data)
}

func (s *Secret) Bytes() []byte {
	return s.data
}

func (s *Secret) Set(data []byte) {
	s.data = data
}

func (a aliasSecret) Text() string {
	return a.S.Text()
}
