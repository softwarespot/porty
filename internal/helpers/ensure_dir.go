package helpers

import (
	"fmt"
	"os"
)

// EnsureDir ensures the directory exists, creating it if it doesn't exist
func EnsureDir(dir string, mode os.FileMode) error {
	if err := os.MkdirAll(dir, mode); err != nil && !os.IsExist(err) {
		return fmt.Errorf("unable to create the directory %q: %w", dir, err)
	}
	return nil
}
