package cmd

import (
	"fmt"

	"github.com/softwarespot/porty/internal/ports"
	"github.com/spf13/cobra"
)

func newRemoveCmd(opts *cliOptions) *cobra.Command {
	return &cobra.Command{
		Use: "remove APPNAME",
		Example: expandExecutable(`Remove a port number for the app name
	$ <EXE> remove myapp`),
		Short:             "Remove the port number for an app name",
		Args:              cobra.MinimumNArgs(1),
		ValidArgsFunction: createCompleteAppNames(opts),
		RunE: func(_ *cobra.Command, args []string) error {
			return runPortsFunc(opts, func(username string, m *ports.Manager) error {
				appName := args[0]
				up, err := m.Unregister(username, appName)
				if err != nil {
					return err
				}

				if opts.flags.asJSON {
					opts.jsonLogger.Log(map[string]any{
						"message": "removed port number",
						"port":    up.Port,
					})
					return nil
				}

				opts.logger.Log(fmt.Sprintf("removed the port number %d for app name %q", up.Port, appName))
				return nil
			})
		},
	}
}
