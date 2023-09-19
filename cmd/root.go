package cmd

import (
	_ "embed"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"

	initialise "github.com/vision-cli/vision/cmd/init"
	"github.com/vision-cli/vision/common/execute"
	cc "github.com/vision-cli/vision/common/plugins"
	rp "github.com/vision-cli/vision/remote-plugins"
)

func init() {
	rootCmd.AddCommand(initialise.RootCmd)
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

//go:embed example.txt
var exampleText string

var rootCmd = &cobra.Command{
	Use:     "vision",
	Short:   "A developer productivity tool",
	Long:    `Vision is tool to create microservice platforms and microservice scaffolding code`,
	Example: exampleText,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Error(err)
	}
}
