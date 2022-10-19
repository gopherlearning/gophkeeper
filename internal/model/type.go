package model

type SecretType uint8

const (
	TextType SecretType = iota
	BinaryType
)

func (d SecretType) String() string {
	return [...]string{"Текс", "Бинарные данные"}[d]
}
