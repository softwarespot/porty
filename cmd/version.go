package cmd

import (
	"math"
	"time"

	"github.com/gosuri/uitable"
	"github.com/softwarespot/porty/internal/helpers"
	"github.com/softwarespot/porty/internal/ports"
	"github.com/softwarespot/porty/internal/version"
	"github.com/spf13/cobra"
)

func newVersionCmd(opts *cliOptions) *cobra.Command {
	return &cobra.Command{
		Use: "version",
		Example: expandExecutable(`Get details about the "<EXE>" version
	$ <EXE> version

Get details about the "<EXE>" version, including "debugging" information
	$ <EXE> version --debug`),
		Short: "Display the application's version",
		RunE: func(_ *cobra.Command, _ []string) error {
			if opts.flags.asJSON {
				out := map[string]any{
					"version":      version.Version,
					"buildTime":    version.Time,
					"buildversion": version.User,
					"goVersion":    version.GoVersion,
				}

				path, err := getDatabasePath(opts)
				if err != nil {
					out["path"] = err.Error()
					opts.jsonLogger.Log(out)
					return nil
				}

				out["path"] = path
				opts.jsonLogger.Log(out)
				return nil
			}

			port := `____________
|  __  __  |
| |  ||  | |
| |__||__| |
|  __  __()|
| |  ||  | |
| |__||__| |
|__________|
`
			opts.logger.Log(port)

			table := uitable.New()
			table.MaxColWidth = 64

			table.AddRow("APP VERSION", version.Version)
			if t, err := time.Parse("2006-01-02T15:04:05-0700", version.Time); err == nil {
				table.AddRow("BUILD TIME", helpers.FormatAsDateTime(t))
			} else {
				table.AddRow("BUILD TIME", version.Time)
			}
			table.AddRow("BUILD USER", version.User)
			table.AddRow("GO VERSION", version.GoVersion)

			opts.logger.Log(table.String())

			table = uitable.New()
			table.MaxColWidth = 64

			path, err := getDatabasePath(opts)
			if err != nil {
				table.AddRow("PORTS FILE", err.Error())
				opts.logger.Log("\n" + table.String())
				return nil
			}

			m, err := ports.Load(path)
			if err != nil {
				table.AddRow("PORTS FILE", err.Error())
				opts.logger.Log("\n" + table.String())
				return nil
			}
			defer m.Close()

			info := m.Info()
			table.AddRow("PORTS FILE", info.Path)

			table.AddRow("MINIMUM PORT", info.MinPort)
			table.AddRow("MAXIMUM PORT", info.MaxPort)

			if port, err := m.Next(); err != nil {
				table.AddRow("NEXT PORT", err.Error())
			} else {
				table.AddRow("NEXT PORT", port)
			}

			ups := m.All(ports.SortByUsernameAppName)
			usernames := map[string]struct{}{}
			for _, up := range ups {
				usernames[up.Username] = struct{}{}
			}

			table.AddRow("USERNAMES COUNT", len(usernames))
			table.AddRow("PORTS COUNT", len(ups))

			avgPerUsername := 0.0
			if len(usernames) > 0 {
				avgPerUsername = math.Ceil(float64(len(ups)) / float64(len(usernames)))
			}
			table.AddRow("PORTS PER USERNAME (AVG)", avgPerUsername)

			opts.logger.Log("\n" + table.String())
			return nil
		},
	}
}
