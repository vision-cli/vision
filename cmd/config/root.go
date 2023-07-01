package config

import (
	"github.com/spf13/cobra"
	"github.com/vision-cli/vision/flag"
)

func init() {
	RootCmd.AddCommand(createCmd)
	createCmd.Flags().AddFlagSet(flag.ConfigFlagset())
}

var RootCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage project config",
	Long:  "Create and modify a project config file",
}
