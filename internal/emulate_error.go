package internal

import (
	"errors"
	"fmt"
	"strings"
)

// emulatedError используется для Эмуляции ошибок
var EmulatedError string

// emulateError используется для Эмуляции ошибок в тесте, для тех функций, в которых невозможно замокать интерфейс
func EmulateError(err error, pos int) error {
	if err != nil {
		return err
	}
	if len(EmulatedError) != 0 && strings.Contains(EmulatedError, fmt.Sprint(pos)) {
		return errors.New(EmulatedError)
	}
	return nil
}
