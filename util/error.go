package util

import (
	"errors"
	"fmt"
)

func E(msg string, args ...any) error {
	return errors.New(fmt.Sprintf(msg, args...))
}

func W(err error, msg string) error {
	return errors.Join(E(msg), err)
}

func J(errs ...error) error {
	return errors.Join(errs...)
}
