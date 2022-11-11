package client

import "github.com/gopherlearning/gophkeeper/internal/model"

type Repository interface {
	Update(m model.Secret) (err error)
	ListKeys() []string
	Get(key string) []byte
}
