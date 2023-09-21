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

// Finds available plugins and initialises them into commands
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
	var arg string
	if len(args) < 1 { // prevents index out of range
		arg = ""
		log.Warnf("No argument provided. Try: \n\n\t vision %v -v", cmd.Use)
	} else {
		arg = args[0]
	}
	exe.RunCommand(cmd.Use, arg)
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

// Seaches all paths on the system to find binaries with specific vision formatting and assigns them to a map.
// The formatting is `vision-plugin-[plugin-name]-[version-number]`
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
