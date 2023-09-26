package cmd

import "github.com/spf13/cobra"

var VersionRootCmd = &cobra.Command{
	Use:   "version",
	Short: "the plugin version",
	Long:  "ditto",
	RunE:  versionCmd,
}

var versionCmd = func(cmd *cobra.Command, args []string) error {

	return nil
}
