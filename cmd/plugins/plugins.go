package plugins

import "github.com/spf13/cobra"

func init() {
	PluginsCmd.AddCommand(listCmd)
}

var PluginsCmd = &cobra.Command{
	Use:   "plugins",
	Short: "Manage plugins",
	Long:  "Manage plugins",
}
