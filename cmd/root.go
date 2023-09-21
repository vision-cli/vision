package cmd

import (
	_ "embed"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"

	initialise "github.com/vision-cli/vision/cmd/init"
	"github.com/vision-cli/vision/execute"
)

func init() {
	rootCmd.AddCommand(initialise.RootCmd)
	pluginMap := findVisionPlugins()
	for plName, plVersion := range pluginMap {
		cobraCmd := &cobra.Command{
			Use:     plName,
			Version: plVersion,
			RunE:    pluginCommand,
		}
		rootCmd.AddCommand(cobraCmd)
	}
}

var pluginCommand = func(cmd *cobra.Command, args []string) error {
	exe := execute.NewPluginExecutor()
	exe.RunCommand(cmd.Use, args[0])
	return nil
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

func findVisionPlugins() map[string]string {
	var pluginMap = make(map[string]string)

	for _, path := range strings.Split(os.Getenv("PATH"), ":") {
		filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
			if info == nil {
				return nil
			}
			if info.IsDir() {
				return nil
			}
			if strings.HasPrefix(info.Name(), "vision-plugin") {
				pluginSplit := strings.Split(info.Name(), "-")

				pluginMap[pluginSplit[2]] = pluginSplit[3]
			}
			return nil
		})
	}
	return pluginMap
}
