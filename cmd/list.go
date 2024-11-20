package cmd

import (
	"fmt"

	"github.com/gosuri/uitable"
	"github.com/mergestat/timediff"
	"github.com/softwarespot/porty/internal/helpers"
	"github.com/softwarespot/porty/internal/ports"
	"github.com/spf13/cobra"
)

type listFlags struct {
	all    bool
	sortBy string
}

func newListCmd(opts *cliOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use: "list",
		Example: expandExecutable(`List all port numbers
	$ <EXE> list

List all port numbers for all users
	$ <EXE> list --all`),
		Short: "List all port numbers",
		RunE: func(_ *cobra.Command, _ []string) error {
			return runPortsFunc(opts, func(username string, m *ports.Manager) error {
				sortBy, err := m.ToSortBy(opts.flags.list.sortBy)
				if err != nil {
					return err
				}

				var ups []ports.UserPort
				if opts.flags.list.all {
					ups = m.All(sortBy)
				} else {
					if ups, err = m.AllByUsername(username, sortBy); err != nil {
						return err
					}
				}

				if opts.flags.asJSON {
					opts.jsonLogger.Log(ups)
					return nil
				}

				table := uitable.New()
				table.MaxColWidth = 64
				table.AddRow("USERNAME", "APPNAME", "PORT", "CREATED AT", "LAST ACCESSED")
				for _, up := range ups {
					table.AddRow(up.Username, up.AppName, up.Port, helpers.FormatAsDateTime(up.CreatedAt), timediff.TimeDiff(up.AccessedAt))
				}

				opts.logger.Log(table.String())
				return nil
			})
		},
	}

	flags := cmd.PersistentFlags()
	flags.BoolVar(&opts.flags.list.all, "all", false, "Enable showing ports for all users")
	flags.StringVar(&opts.flags.list.sortBy, "sort-by", ports.SortByUsernameAppName.String(), fmt.Sprintf("Sort the ports by. Must be one of: %s", helpers.ShellQuoteJoin(ports.SortByStrings)))

	cmd.RegisterFlagCompletionFunc("sort-by", func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
		return ports.SortByStrings, cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}
