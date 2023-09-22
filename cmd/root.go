package cmd

import (
	_ "embed"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	initialise "github.com/vision-cli/vision/cmd/init"
	"github.com/vision-cli/vision/internal/plugin"
)

// Finds available plugins and initialises them into commands
func init() {
	rootCmd.AddCommand(initialise.RootCmd)
	rootCmd.Flags().AddFlagSet(initVisionFlags())
	plugins := findVisionPlugins()
	for _, plugin := range plugins {
		cmd, err := createCommand(plugin)
		if err != nil {
			// TODO(steve): handle broken commands, maybe a vision doctor command???
			log.Error(plugin.FullPath, err)
			continue
		}
		rootCmd.AddCommand(cmd)
	}
}

// createCommand takes in a plugin and returns a cobra command to interact with that plugin
func createCommand(p pluginPath) (*cobra.Command, error) {
	exe := plugin.NewExecutor(p.FullPath)
	info, err := exe.Info()
	if err != nil {
		return nil, err
	}
	version, err := exe.Version()
	if err != nil {
		return nil, err
	}

	cobraCmd := &cobra.Command{
		Use:     p.Name,
		Version: version.SemVer,
		Short:   info.ShortDescription,
		Long:    info.LongDescription,
		RunE:    createPluginCommandHandler(p),
	}

	return cobraCmd, nil
}

func createPluginCommandHandler(p pluginPath) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 { // prevents index out of range
			log.Warnf("No argument provided. Try: \n\t\n vision %v -v", cmd.Use)
		}
		exe := plugin.NewExecutor(p.FullPath)
		switch args[0] {
		case "init":
			i, err := exe.Init()
			if err != nil {
				return err
			}
			// TODO merge into vison config
			log.Info(i.Config)
		}
		return nil
	}
}

func initVisionFlags() *pflag.FlagSet {
	fs := pflag.NewFlagSet("vision", 1)
	return fs
}

type pluginPath struct {
	Name     string
	Version  string
	FullPath string
}

// Seaches all dirs in the PATH envar to find binaries with specific vision formatting and assigns them to a map.
// The formatting is `vision-plugin-[plugin-name]-[version-number]`
func findVisionPlugins() []pluginPath {
	const prefix = "vision-plugin"

	var plugins []pluginPath

	sysPath := os.Getenv("PATH")
	paths := strings.Split(sysPath, ":")

	m := make(map[string]struct{})
	for _, p := range paths {
		m[p] = struct{}{}
	}
	var uniqPaths []string
	for k := range m {
		uniqPaths = append(uniqPaths, k)
	}
	paths = uniqPaths

	for _, path := range uniqPaths {
		filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
			if info == nil {
				return nil
			}
			if info.IsDir() {
				return nil
			}
			if strings.HasPrefix(info.Name(), prefix) {
				pluginSplit := strings.Split(info.Name(), "-")
				if len(pluginSplit) == 4 { // only process correctly named plugins
					name, version := pluginSplit[2], pluginSplit[3]
					plugins = append(plugins, pluginPath{
						Name:     name,
						Version:  version,
						FullPath: path,
					})
				}
			}
			return nil
		})
	}
	return plugins
}

//go:embed example.txt
var exampleText string

var rootCmd = &cobra.Command{
	Use:     "vision",
	Short:   "A developer productivity tool",
	Long:    `Vision is a tool to create microservice platforms and microservice scaffolding code`,
	Example: exampleText,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Error(err)
	}
}
