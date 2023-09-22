package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"

	"github.com/vision-cli/vision/internal/plugin"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "doctor",
	Short: "check status of plugins",
	Long:  "Check the status of all plugins and their commands. If a command fails, it will log an error.",
	RunE:  doctorCommand,
}

var healthRecord []string

var doctorCommand = func(cmd *cobra.Command, args []string) error {
	plugins := FindVisionPlugins()

	for _, plug := range plugins {
		log.Infof("The plugins available are: %v", plug.Name)
	}

	for _, plug := range plugins {
		// call each of the built in commands
		exe := plugin.NewExecutor(plug.FullPath)
		info, err := exe.Info()
		if err != nil {
			healthCheck := fmt.Sprintf("Info is faulty for plugin: %v", plug.Name)
			healthRecord = append(healthRecord, healthCheck)
		} else if info.ShortDescription == "" || info.LongDescription == "" {
			healthCheck := fmt.Sprintf("Info is faulty for plugin: %v", plug.Name)
			healthRecord = append(healthRecord, healthCheck)
		}

		ini, err := exe.Init()
		if err != nil {
			healthCheck := fmt.Sprintf("Init is faulty for plugin: %v", plug.Name)
			healthRecord = append(healthRecord, healthCheck)
		} else if ini.Config == "" || ini.Config == nil {
			healthCheck := fmt.Sprintf("Init is faulty for plugin: %v", plug.Name)
			healthRecord = append(healthRecord, healthCheck)
		}

		vers, err := exe.Version()
		if err != nil {
			healthCheck := fmt.Sprintf("Version is faulty for plugin: %v", plug.Name)
			healthRecord = append(healthRecord, healthCheck)
		} else if vers.SemVer == "" {
			healthCheck := fmt.Sprintf("Version is faulty for plugin: %v", plug.Name)
			healthRecord = append(healthRecord, healthCheck)
		}
	}
	for _, hr := range healthRecord {
		log.Warn(hr)
	}
	return nil
}

type PluginPath struct {
	Name     string
	Version  string
	FullPath string
}

// Searches all dirs in the PATH envar to find binaries with specific vision formatting and assigns them to a map.
// The formatting is `vision-plugin-[plugin-name]-[version-number]`
func FindVisionPlugins() []PluginPath {
	const prefix = "vision-plugin"

	var plugins []PluginPath

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
					plugins = append(plugins, PluginPath{
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
