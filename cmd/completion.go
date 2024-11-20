package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/softwarespot/porty/internal/helpers"
	"github.com/spf13/cobra"
)

func newCompletionCmd() *cobra.Command {
	byShellFuncs := map[string]func(cmd *cobra.Command, out io.Writer) error{
		"bash": func(cmd *cobra.Command, w io.Writer) error {
			return cmd.Root().GenBashCompletion(w)
		},
		"fish": func(cmd *cobra.Command, w io.Writer) error {
			return cmd.Root().GenFishCompletion(w, true)
		},
		"zsh": func(cmd *cobra.Command, w io.Writer) error {
			return cmd.Root().GenZshCompletion(w)
		},
	}

	var shells []string
	for s := range byShellFuncs {
		shells = append(shells, s)
	}

	shellList := helpers.ShellQuoteJoin(shells)

	return &cobra.Command{
		Use:   "completion SHELL",
		Short: fmt.Sprintf("Generate autocompletions script for the specified shell (%s)", shellList),
		Long: expandExecutable(`Generate autocompletions script for the specified shells (%s)

This command generates shell autocompletions for bash e.g.
	$ <EXE> completion bash

Can be sourced as such in .bashrc
	$ source <(<EXE> completion bash)`, shellList),
		ValidArgs: shells,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("shell must be defined, expected one of the following: %s", shellList)
			}

			shell := args[0]
			fn, ok := byShellFuncs[shell]
			if !ok {
				return fmt.Errorf("invalid shell of %q was provided, expected one of the following: %s", shell, shellList)
			}
			return fn(cmd, os.Stdout)
		},
	}
}
