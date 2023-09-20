package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/vision-cli/vision/execute"
)

func init() {
	RootCmd.Flags().AddFlagSet(pluginFlags())
}

var (
	getVersion bool
)

func pluginFlags() *pflag.FlagSet {
	fs := pflag.NewFlagSet("project", 1)
	fs.BoolVarP(&getVersion, "version", "v", false, "return the version")
	return fs
}

var RootCmd = &cobra.Command{
	Use:   "project",
	Short: "Create a new project",
	Long:  "Ditto",
	RunE: func(cmd *cobra.Command, args []string) error {
		if getVersion {
			// TODO (luke): Implement real plugins
			exe := execute.NewPluginExecutor("/opt/homebrew/bin/go")
			if err := exe.Version(); err != nil {
				return err
			}
		}

		return nil
	},
}
