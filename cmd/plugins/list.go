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
	RunE:    cmd,
}

var cmd = func(cmd *cobra.Command, args []string) error {
	plugins, err := plugin.Find()
	if err != nil {
		return err
	}
	for _, p := range plugins {
		fmt.Println(p.Name)
	}
	return nil
}
