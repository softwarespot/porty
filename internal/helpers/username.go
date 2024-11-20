package helpers

import (
	"fmt"
	"os/user"
)

// Username returns the current username
func Username() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("unable to get the current username: %w", err)
	}
	return u.Username, nil
}
