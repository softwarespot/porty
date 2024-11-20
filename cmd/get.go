package cmd

import (
	"github.com/softwarespot/porty/internal/ports"
	"github.com/spf13/cobra"
)

func newGetCmd(opts *cliOptions) *cobra.Command {
	return &cobra.Command{
		Use: "get APPNAME",
		Example: expandExecutable(`Get or assign the next available port number for the app name
	$ <EXE> get myapp`),
		Short:             "Get the port number for an app name; otherwise assigns a new port number",
		Args:              cobra.MinimumNArgs(1),
		ValidArgsFunction: createCompleteAppNames(opts),
		RunE: func(_ *cobra.Command, args []string) error {
			return runPortsFunc(opts, func(username string, m *ports.Manager) error {
				appName := args[0]
				if up, err := m.Get(username, appName); err == nil {
					if opts.flags.asJSON {
						opts.jsonLogger.Log(map[string]any{
							"message": "got port number",
							"port":    up.Port,
						})
					} else {
						opts.logger.Log(up.Port.String())
					}
					return nil
				}

				up, err := m.Register(username, appName)
				if err != nil {
					return err
				}

				if opts.flags.asJSON {
					opts.jsonLogger.Log(map[string]any{
						"message": "got port number",
						"port":    up.Port,
					})
					return nil
				}

				opts.logger.Log(up.Port.String())
				return nil
			})
		},
	}
}
