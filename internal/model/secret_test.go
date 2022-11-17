package model

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testsTestSecretRead = []struct {
	secret *Secret
	name   string
	text   string
	result string
}{
	{
		name: "success text",
		text: "hello",
		result: `---
hello
---
labels:
  bla: bla
name: success-text
owner: testovich
`,
		secret: &Secret{
			Name:   "success-text",
			Owner:  "testovich",
			Labels: map[string]string{"bla": "bla"},
			Type:   TextType,
			data:   []byte("hello"),
		},
	},
	{
		name: "success text without labels",
		text: "hello",
		result: `---
hello
---
name: success-text
owner: testovich
`,
		secret: &Secret{
			Name:  "success-text",
			Owner: "testovich",
			Type:  TextType,
			data:  []byte("hello"),
		},
	},
	{
		name:   "success binary",
		text:   "это не текстовые данные. Тип - Бинарные данные",
		result: "---\nэто не текстовые данные. Тип - Бинарные данные\n---\nname: success-text\nowner: testovich\n",
		secret: &Secret{
			Name:  "success-text",
			Owner: "testovich",
			Type:  BinaryType,
			data:  []byte("hello"),
		},
	},
}

func TestSecretRead(t *testing.T) {
	for _, v := range testsTestSecretRead {
		t.Run(v.name, func(t *testing.T) {
			assert.Equal(t, v.secret.Text(), v.text)
			assert.Equal(t, v.secret.String(), v.result)
			b := v.secret.Bytes()
			assert.NotEmpty(t, b)
			assert.Equal(t, b, v.secret.data)
		})
	}
}

func TestSecretSet(t *testing.T) {
	var testsTestSecretRead = []struct {
		secret *Secret
		// cryptor *FakeDecryptorEncryptor
		name  string
		bytes []byte
	}{
		{
			name:   "success",
			secret: &Secret{},
			// cryptor: &FakeDecryptorEncryptor{},
			bytes: []byte("success"),
		},
		{
			name:   "error",
			secret: &Secret{},
			// cryptor: &FakeDecryptorEncryptor{Err: errors.New("fake error")},
			bytes: nil,
		},
	}

	for _, v := range testsTestSecretRead {
		t.Run(v.name, func(t *testing.T) {
			fmt.Println(v.secret.Name)
			assert.Empty(t, v.secret.Bytes())
			v.secret.Set(v.bytes)
			if len(v.bytes) != 0 {
				assert.NotEmpty(t, v.secret.Bytes())
			} else {
				assert.Empty(t, v.secret.Bytes())
			}
		})
	}
}

type FakeDecryptorEncryptor struct {
	Err error
}

func (f *FakeDecryptorEncryptor) error(b []byte) ([]byte, error) {
	if f.Err != nil {
		return nil, f.Err
	}

	return b, nil
}

func (f *FakeDecryptorEncryptor) Decrypt(b []byte) ([]byte, error) {
	return f.error(b)
}

func (f *FakeDecryptorEncryptor) Encrypt(b []byte, recipient ...string) ([]byte, error) {
	return f.error(b)
}
