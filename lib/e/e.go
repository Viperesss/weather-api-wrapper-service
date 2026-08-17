package e

import "fmt"

func Wrap(text string, err error) error {
	return fmt.Errorf("%s: %w", text, err)
}
