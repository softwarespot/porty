package helpers

import (
	"os"
	"path/filepath"
)

// ExecutableName returns the executable's name
func ExecutableName() string {
	return filepath.Base(os.Args[0])
}
