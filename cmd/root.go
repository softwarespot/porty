package cmd

import (
	"os"

	"github.com/softwarespot/porty/internal/env"
	"github.com/softwarespot/porty/internal/helpers"
	"github.com/softwarespot/porty/internal/logging"
	"github.com/spf13/cobra"
)

type cliFlags struct {
	dir string

	list listFlags
	who  whoFlags

	asJSON bool
	help   bool
}

type cliOptions struct {
	flags  cliFlags
	dbName string

	logger     *logging.Logger
	jsonLogger *logging.JSONLogger
}

func newRootCmd(opts *cliOptions) (*cobra.Command, error) {
	// See URL: https://qua.name/antolius/making-a-testable-cobra-cli-app
	cmd := &cobra.Command{
		Use: helpers.ExecutableName(),
		Example: expandExecutable(`Get details about the "<EXE>" version
	$ <EXE> version

Get or assign the next available port number for the app name
	$ <EXE> get myapp

Remove a port number for the app name
	$ <EXE> remove myapp

List all port numbers
	$ <EXE> list

Get the next available port number, without assigning to an app name
	$ <EXE> next

Get who has been assigned to a port number
	$ <EXE> who 8000

Get details about the "<EXE>" version, including "debugging" information
	$ <EXE> version

Example of adding autocompletion for "bash" to the ".bashrc" file
	$ echo 'source <(<EXE> completion bash)' >> ~/.bashrc
	$ source ~/.bashrc`),
		Short: "Manage port numbers for the current user",
		Long: `
Manage port numbers for the current user

See https://github.com/softwarespot/porty for more details about the application
`,
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			// Workaround for an issue where the "flags" are parsed only after this function
			// has been called
			if cmd.Name() == cobra.ShellCompRequestCmd || cmd.Name() == cobra.ShellCompNoDescRequestCmd {
				return nil
			}
			if opts.flags.asJSON {
				opts.jsonLogger = logging.NewJSONLogger()
			}
			return nil
		},
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.SetHelpCommand(&cobra.Command{
		Hidden: true,
	})

	cmd.AddCommand(
		newCompletionCmd(),
		newGetCmd(opts),
		newInitCmd(opts),
		newListCmd(opts),
		newNextCmd(opts),
		newRemoveCmd(opts),
		newWhoCmd(opts),
		newVersionCmd(opts),
	)

	dir, err := getPortyDir()
	if err != nil {
		return nil, err
	}

	flags := cmd.PersistentFlags()
	flags.StringVar(&opts.flags.dir, "dir", env.Get("PORTY_DIR", dir), "Ports database directory")
	flags.BoolVarP(&opts.flags.asJSON, "json", "j", false, "Output the result as JSON")
	flags.BoolVarP(&opts.flags.help, "help", "h", false, "Display the application's usage")

	return cmd, nil
}

// Execute adds all child commands to the root command and sets flags appropriately
func Execute() {
	opts := &cliOptions{
		flags:  cliFlags{},
		dbName: "porty_db.json",
		logger: logging.NewLogger(),
	}

	root, err := newRootCmd(opts)
	if err != nil {
		opts.logger.LogError(err)
		os.Exit(1)
		return
	}

	if err := root.Execute(); err != nil {
		opts.logger.LogError(err)
		os.Exit(1)
	}
}
