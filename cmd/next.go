package cmd

import (
	"fmt"

	"github.com/softwarespot/porty/internal/ports"
	"github.com/spf13/cobra"
)

func newNextCmd(opts *cliOptions) *cobra.Command {
	return &cobra.Command{
		Use: "next",
		Example: expandExecutable(`Get the next available port number, without assigning to an app name
	$ <EXE> next`),
		Short: "Get the next available port number, without assigning to an app name",
		RunE: func(_ *cobra.Command, _ []string) error {
			return runPortsFunc(opts, func(_ string, m *ports.Manager) error {
				port, err := m.Next()
				if err != nil {
					return err
				}

				if opts.flags.asJSON {
					opts.jsonLogger.Log(map[string]any{
						"message": "next port number",
						"port":    port,
					})
					return nil
				}

				opts.logger.Log(fmt.Sprintf("next port number is %d", port))
				return nil
			})
		},
	}
}
