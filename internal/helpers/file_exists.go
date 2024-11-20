package helpers

import "os"

// FileExists determines if a filepath exists and is not a directory
func FileExists(filepath string) bool {
	f, err := os.Stat(filepath)
	return err == nil && !f.IsDir()
}
