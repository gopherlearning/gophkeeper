package client

import "github.com/gopherlearning/gophkeeper/internal/model"

type Repository interface {
	Update(m model.Secret) (err error)
	ListKeys(...model.SecretType) []string
	Get(m model.Secret) *model.Secret
	Remove(m model.Secret) error
}
