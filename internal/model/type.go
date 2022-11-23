package model

type SecretType uint8

const (
	TextType SecretType = iota
	BinaryType
	CardType
	PasswordType
)

func (d SecretType) String() string {
	return [...]string{"Текст", "Бинарные данные", "Банковская карта", "Логин/Пароль"}[d]
}
