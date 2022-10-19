package model

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testsTestSecretRead = []struct {
	secret  *Secret
	cryptor *FakeDecryptorEncryptor
	name    string
	text    string
	result  string
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
			Data:   []byte("hello"),
		},
		cryptor: &FakeDecryptorEncryptor{},
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
			Data:  []byte("hello"),
		},
		cryptor: &FakeDecryptorEncryptor{},
	},
	{
		name:   "success binary",
		text:   "это не текстовые данные. Тип - Бинарные данные",
		result: "---\nэто не текстовые данные. Тип - Бинарные данные\n---\nname: success-text\nowner: testovich\n",
		secret: &Secret{
			Name:  "success-text",
			Owner: "testovich",
			Type:  BinaryType,
			Data:  []byte("hello"),
		},
		cryptor: &FakeDecryptorEncryptor{},
	},
	{
		name:    "decryptor error",
		text:    "не удалось расшифровать данные: fake error",
		result:  "---\nне удалось расшифровать данные: fake error\n---\nname: \nowner: \n",
		secret:  &Secret{},
		cryptor: &FakeDecryptorEncryptor{Err: errors.New("fake error")},
	},
}

func TestSecretRead(t *testing.T) {
	for _, v := range testsTestSecretRead {
		t.Run(v.name, func(t *testing.T) {
			assert.Equal(t, v.secret.Text(v.cryptor), v.text)
			assert.Equal(t, v.secret.String(v.cryptor), v.result)
			b, err := v.secret.Bytes(v.cryptor)
			if v.cryptor.Err != nil {
				assert.ErrorContains(t, err, v.cryptor.Err.Error())
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, b, v.secret.Data)
		})
	}
}

func TestSecretSet(t *testing.T) {
	var testsTestSecretRead = []struct {
		secret  *Secret
		cryptor *FakeDecryptorEncryptor
		name    string
		bytes   []byte
	}{
		{
			name:    "success",
			secret:  &Secret{},
			cryptor: &FakeDecryptorEncryptor{},
			bytes:   []byte("success"),
		},
		{
			name:    "error",
			secret:  &Secret{},
			cryptor: &FakeDecryptorEncryptor{Err: errors.New("fake error")},
			bytes:   nil,
		},
	}

	for _, v := range testsTestSecretRead {
		t.Run(v.name, func(t *testing.T) {
			err := v.secret.Set(v.cryptor, v.bytes)
			if v.cryptor.Err != nil {
				assert.ErrorContains(t, err, v.cryptor.Err.Error())
				return
			}
			assert.NoError(t, err)
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

func (f *FakeDecryptorEncryptor) Encrypt(b []byte) ([]byte, error) {
	return f.error(b)
}
