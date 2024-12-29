package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/softwarespot/porty/internal/helpers"
	"github.com/softwarespot/porty/internal/ports"
)

func expandExecutable(format string, a ...any) string {
	return strings.ReplaceAll(fmt.Sprintf(format, a...), "<EXE>", helpers.ExecutableName())
}

func getDatabaseDir() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("unable to get the user configuration directory")
	}
	return filepath.Join(dir, "porty"), nil
}

func getDatabasePath(opts *cliOptions) (string, error) {
	path := filepath.Join(opts.flags.dir, opts.dbName)
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("unable to get the absolute directory of the ports database %q: %w", path, err)
	}
	return absPath, nil
}

func execPortsFunc(opts *cliOptions, fn func(username string, m *ports.Manager) error) error {
	username, err := helpers.Username()
	if err != nil {
		return err
	}

	path, err := getDatabasePath(opts)
	if err != nil {
		return err
	}

	if !helpers.FileExists(path) {
		return fmt.Errorf(`ports database %q doesn't exist. Run "%s init" to initialize the ports database`, path, helpers.ExecutableName())
	}

	m, err := ports.Load(path)
	if err != nil {
		return err
	}

	// Ignore the error
	defer m.Close()

	if err := fn(username, m); err != nil {
		return err
	}
	return nil
}
