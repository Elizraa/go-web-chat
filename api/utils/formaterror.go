package utils

import (
	"errors"
	"strings"
)

// Cleanly formatted error response
func FormatError(err string) error {

	if strings.Contains(err, "url") {
		return errors.New("Missing URL")
	}

	return errors.New(err)
}
