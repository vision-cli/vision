package cmd

import (
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"

	"my-test-module/cmd/generate"
	"my-test-module/cmd/info"
	"my-test-module/cmd/initialise"
	"my-test-module/cmd/version"
)

func init() {
	rootCmd.AddCommand(initialise.InitCmd)
	rootCmd.AddCommand(info.InfoCmd)
	rootCmd.AddCommand(version.VersionCmd)
	rootCmd.AddCommand(generate.GenerateCmd)
}

var rootCmd = &cobra.Command{
	Use:                "plugin",
	FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Error(err)
	}
}
