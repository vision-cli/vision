package cmd

import (
	_ "embed"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"

	initialise "github.com/vision-cli/vision/cmd/init"
	plug "github.com/vision-cli/vision/cmd/plugin"
	"github.com/vision-cli/vision/common/execute"
	cc "github.com/vision-cli/vision/common/plugins"
	rp "github.com/vision-cli/vision/remote-plugins"
)

func init() {
	rootCmd.AddCommand(initialise.RootCmd)
	rootCmd.AddCommand(plug.RootCmd)
	osExecutor := execute.NewOsExecutor()
	p, err := cc.GetPlugins(osExecutor)
	if err != nil {
		log.Warn("cannot get plugins: %v", err)
	}
	for _, pl := range p {
		cobraCmd, err := rp.GetCobraCommand(pl, osExecutor)
		if err != nil {
			log.Warn("cannot get cobra command %s: %v", pl.Name, err)
		}
		rootCmd.AddCommand(cobraCmd)
	}
}

//go:embed vision-help.txt
var visionHelp string

var rootCmd = &cobra.Command{
	Use:     "vision",
	Short:   "A developer productivity tool",
	Long:    `Vision is tool to create microservice platforms and microservice scaffolding code`,
	Example: visionHelp,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Error(err)
	}
}
