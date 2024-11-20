package cmd

import (
	"fmt"

	"github.com/softwarespot/porty/internal/helpers"
	"github.com/softwarespot/porty/internal/ports"

	"github.com/spf13/cobra"
)

func newInitCmd(opts *cliOptions) *cobra.Command {
	return &cobra.Command{
		Use: "init",
		Example: expandExecutable(`Initialize the directory with the ports database
	$ <EXE> init`),
		Short: "Initialize the directory with the ports database",
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := helpers.EnsureDir(opts.flags.dir, 0o700); err != nil {
				return err
			}

			path, err := getDatabasePath(opts)
			if err != nil {
				return err
			}

			if err := ports.Init(path); err != nil {
				return err
			}

			if opts.flags.asJSON {
				opts.jsonLogger.Log(map[string]any{
					"message": "created ports database",
					"path":    path,
				})
				return nil
			}

			opts.logger.Log(fmt.Sprintf("created ports database %q", path))
			return nil
		},
	}
}
