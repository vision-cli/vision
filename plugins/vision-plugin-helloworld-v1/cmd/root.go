package cmd

import (
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"

	info "github.com/vision-cli/vision/plugins/vision-plugin-helloworld-v1/cmd/info"
	initialise "github.com/vision-cli/vision/plugins/vision-plugin-helloworld-v1/cmd/initialise"
	version "github.com/vision-cli/vision/plugins/vision-plugin-helloworld-v1/cmd/version"
)

func init() {
	rootCmd.AddCommand(initialise.InitRootCmd)
	rootCmd.AddCommand(info.InfoRootCmd)
	rootCmd.AddCommand(version.VersionRootCmd)
}

var rootCmd = &cobra.Command{
	Use:                "helloworld",
	FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Error(err)
	}
}
