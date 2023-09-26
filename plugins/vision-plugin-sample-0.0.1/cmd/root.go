package cmd

import (
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	info "github.com/vision-cli/vision/plugins/vision-plugin-sample-0.0.1/cmd/info"
)

func init() {
	rootCmd.AddCommand(info.RootCmd)
}

var rootCmd = &cobra.Command{
	Use:   "sample",
	Short: "The sample plugin",
	Long:  "The sample plugin of the vision project",
	RunE:  sampleCmd,
}

var sampleCmd = func(cmd *cobra.Command, args []string) error {

	return nil
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Error(err)
	}
}
