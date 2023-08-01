//go:build !prod

package gojason

import "fmt"

func log(args ...any) {
	fmt.Println(args...)
}
