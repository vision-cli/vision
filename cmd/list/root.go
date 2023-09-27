package list

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/vision-cli/vision/internal/plugin"
)

var RootCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all installed vision plugins.",
	Long:    "",
	Run:     listCmd,
}

var listCmd = func(cmd *cobra.Command, args []string) {
	plugins := plugin.Find()
	for _, p := range plugins {
		fmt.Println(p.Name)
	}
}
