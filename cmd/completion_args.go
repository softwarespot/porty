package cmd

import (
	"github.com/softwarespot/porty/internal/ports"
	"github.com/spf13/cobra"
)

func createCompleteAppNames(opts *cliOptions) func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(_ *cobra.Command, args []string, _ string) ([]string, cobra.ShellCompDirective) {
		if len(args) > 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		var appnames []string

		// Ignore the error
		runPortsFunc(opts, func(username string, m *ports.Manager) error {
			ups, err := m.AllByUsername(username, ports.SortByUsernameAppName)
			if err != nil {
				return err
			}
			for _, up := range ups {
				appnames = append(appnames, up.AppName)
			}
			return nil
		})
		return appnames, cobra.ShellCompDirectiveNoFileComp
	}
}

func createCompletePorts(opts *cliOptions) func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(_ *cobra.Command, args []string, _ string) ([]string, cobra.ShellCompDirective) {
		if len(args) > 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		var ps []string

		// Ignore the error
		runPortsFunc(opts, func(_ string, m *ports.Manager) error {
			for _, up := range m.All(ports.SortByUsernameAppName) {
				ps = append(ps, up.Port.String())
			}
			return nil
		})
		return ps, cobra.ShellCompDirectiveNoFileComp
	}
}
