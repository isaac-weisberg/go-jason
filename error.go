package gojason

import (
	"errors"
	"fmt"
)

func e(msg string, args ...any) error {
	return errors.New(fmt.Sprintf(msg, args...))
}

func w(err error, msg string) error {
	return errors.Join(e(msg), err)
}

func j(errs ...error) error {
	return errors.Join(errs...)
}
