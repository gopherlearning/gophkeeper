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

type Decryptor interface {
	Decrypt([]byte) ([]byte, error)
}
type Encryptor interface {
	Encrypt([]byte) ([]byte, error)
}

type Secret struct {
	Name   string
	Owner  string
	Labels map[string]string
	Data   []byte
	Type   SecretType
}

type aliasSecret struct {
	S *Secret
	d Decryptor
	// e Encryptor
}

func (s *Secret) String(d Decryptor) string {
	out := &bytes.Buffer{}
	_ = secretTemplate.Execute(out, aliasSecret{S: s, d: d})

	return out.String()
}

func (s *Secret) Text(d Decryptor) string {
	if s.Type != TextType {
		return fmt.Errorf("это не текстовые данные. Тип - %s", s.Type).Error()
	}

	r, err := d.Decrypt(s.Data)
	if err != nil {
		return fmt.Errorf("не удалось расшифровать данные: %v", err).Error()
	}

	return string(r)
}

func (s *Secret) Bytes(d Decryptor) ([]byte, error) {
	r, err := d.Decrypt(s.Data)
	if err != nil {
		return nil, fmt.Errorf("не удалось расшифровать данные: %v", err)
	}

	return r, nil
}

func (s *Secret) Set(e Encryptor, data []byte) error {
	r, err := e.Encrypt(s.Data)
	if err != nil {
		return fmt.Errorf("не удалось зашифровать данные: %v", err)
	}

	s.Data = r

	return nil
}

func (a aliasSecret) Text() string {
	return a.S.Text(a.d)
}
