package plugins

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/vision-cli/vision/internal/plugin"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all installed vision plugins",
	Long:    "List all installed vision plugins",
	Run:     cmd,
}

var cmd = func(cmd *cobra.Command, args []string) {
	plugins := plugin.Find()
	for _, p := range plugins {
		fmt.Println(p.Name)
	}
}
