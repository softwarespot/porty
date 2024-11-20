package cmd

import (
	"fmt"

	"github.com/gosuri/uitable"
	"github.com/mergestat/timediff"
	"github.com/softwarespot/porty/internal/helpers"
	"github.com/softwarespot/porty/internal/ports"
	"github.com/spf13/cobra"
)

type whoFlags struct {
	sortBy string
}

func newWhoCmd(opts *cliOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use: "who PORT",
		Example: expandExecutable(`Get who has been assigned to a port number
	$ <EXE> who 8000`),
		Short:             "Who gets who has been assigned to a port number",
		Args:              cobra.MinimumNArgs(1),
		ValidArgsFunction: createCompletePorts(opts),
		RunE: func(_ *cobra.Command, args []string) error {
			return runPortsFunc(opts, func(_ string, m *ports.Manager) error {
				port, err := m.ToPort(args[0])
				if err != nil {
					return err
				}

				up, err := m.GetByPort(port)
				if err != nil {
					return err
				}

				if opts.flags.asJSON {
					opts.jsonLogger.Log(up)
					return nil
				}

				table := uitable.New()
				table.MaxColWidth = 64
				table.AddRow("USERNAME", "APPNAME", "PORT", "CREATED AT", "LAST ACCESSED")
				table.AddRow(up.Username, up.AppName, up.Port, helpers.FormatAsDateTime(up.CreatedAt), timediff.TimeDiff(up.AccessedAt))

				opts.logger.Log(table.String())
				return nil
			})
		},
	}

	flags := cmd.PersistentFlags()
	flags.StringVar(&opts.flags.who.sortBy, "sort-by", ports.SortByUsernameAppName.String(), fmt.Sprintf("Sort the ports by. Must be one of: %s", helpers.ShellQuoteJoin(ports.SortByStrings)))

	cmd.RegisterFlagCompletionFunc("sort-by", func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
		return ports.SortByStrings, cobra.ShellCompDirectiveNoFileComp
	})
	return cmd
}
