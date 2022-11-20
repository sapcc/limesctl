package util

import "fmt"

func WrapError(err error, msg string) error {
	return fmt.Errorf("%s: %w", msg, err)
}
