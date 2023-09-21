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
)

func init() {
	rootCmd.AddCommand(initialise.RootCmd)
	pluginArray := findVisionPlugins()
	for _, pl := range pluginArray {
		cobraCmd := &cobra.Command{Use: pl}
		rootCmd.AddCommand(cobraCmd)
	}
}

//go:embed example.txt
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

func findVisionPlugins() []string {
	var pluginNameArray []string

	for _, path := range strings.Split(os.Getenv("PATH"), ":") {
		filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
			if info == nil {
				return nil
			}
			if info.IsDir() {
				return nil
			}
			if strings.HasPrefix(info.Name(), "vision-plugin") {
				nameArray := strings.Split(info.Name(), "-")
				pluginNameArray = append(pluginNameArray, nameArray[2])
			}
			return nil
		})
	}

	return pluginNameArray
}
