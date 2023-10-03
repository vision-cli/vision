package cmd

import (
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"

	"{{.ModuleName}}/cmd/generate"
	"{{.ModuleName}}/cmd/info"
	"{{.ModuleName}}/cmd/initialise"
	"{{.ModuleName}}/cmd/version"
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
