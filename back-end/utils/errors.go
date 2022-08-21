package utils

import "fmt"

func NewErrorWithInnerError(msg string, innerErr error) error {
	return fmt.Errorf(msg+": %w", innerErr)
}
