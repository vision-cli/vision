package plugins

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	api_v1 "github.com/vision-cli/api/v1"
	"github.com/vision-cli/vision/cli"
	"github.com/vision-cli/vision/config"
	"github.com/vision-cli/vision/flag"
	"github.com/vision-cli/vision/placeholders"
)

var UsageQuery = api_v1.PluginRequest{
	Command: api_v1.CommandUsage,
}

func GetPlugins() []string {
	var plugins []string
	pluginPath := goBinPath()
	pluginFiles, err := os.ReadDir(pluginPath)
	if err != nil {
		cli.Fatalf("Cannot read plugin directory %s", pluginPath)
	}
	for _, pluginFile := range pluginFiles {
		if !pluginFile.IsDir() && fileIsVisionPlugin(pluginFile.Name()) {
			plugins = append(plugins, pluginFile.Name())
		}
	}
	return plugins
}

func GetCobraCommand(plugin string) *cobra.Command {
	pluginPath := filepath.Join(goBinPath(), plugin)
	usageQuery, err := Marshal(UsageQuery)
	if err != nil {
		cli.Fatalf("Cannot marshal usage query for plugin %s: %s", plugin, err.Error())
	}
	response := callPlugin(pluginPath, usageQuery)
	usage, err := Unmarshal[api_v1.PluginUsageResponse](response)
	if err != nil {
		cli.Fatalf("Cannot marshal usage query for plugin %s: %s", plugin, err.Error())
	}
	cc := cobra.Command{
		Use:     usage.Use,
		Short:   usage.Short,
		Long:    usage.Long,
		Example: usage.Example,
		Run: func(cmd *cobra.Command, args []string) {
			if err := initializeConfig(cmd); err != nil {
				cli.Fatalf("Cannot initialize config: %v", err)
			}
			placeholders := placeholders.NewPlaceholders(cmd.Flags(), ".", "default", "", "")
			runQuery, err := Marshal(api_v1.PluginRequest{
				Command:      api_v1.CommandRun,
				Args:         args,
				Flags:        []api_v1.PluginFlag{},
				Placeholders: placeholders,
			})
			if err != nil {
				cli.Fatalf("Cannot marshal run query for plugin %s", plugin)
			}
			response := callPlugin(pluginPath, runQuery)
			result, err := Unmarshal[api_v1.PluginResponse](response)
			if err != nil {
				cli.Fatalf("Cannot unmarshal response from plugin %s", plugin)
			}
			if result.Error != "" {
				cli.Fatalf(result.Error)
			}
			cli.Infof(result.Result)
		},
	}
	cc.Flags().AddFlagSet(flag.ConfigFlagset())
	return &cc
}

func Unmarshal[T any](reqStr string) (T, error) {
	var req T
	err := json.Unmarshal([]byte(reqStr), &req)
	if err != nil {
		return req, err
	}
	return req, nil
}

func Marshal[T any](resp T) (string, error) {
	respStr, err := json.Marshal(resp)
	if err != nil {
		return "", err
	}
	return string(respStr), nil
}

func goBinPath() string {
	goBinPath := os.Getenv("GOBIN")
	if goBinPath == "" {
		goPath, err := exec.Command("go", "env", "GOPATH").Output()
		if err != nil {
			cli.Fatalf("Cannot determine GOBIN. Please set GOBIN or GOPATH")
		}
		goBinPath = string(goPath)[:len(goPath)-1] + "/bin"
	}
	return goBinPath
}

func fileIsVisionPlugin(filename string) bool {
	c := strings.Split(filename, "-")
	// eventually remove 'example'
	if len(c) != 4 || c[0] != "vision" || c[1] != "plugin" {
		return false
	}
	return true
}

func callPlugin(plugin string, input string) string {
	cmd := exec.Command(plugin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.Output()
	if err != nil {
		cli.Fatalf("Cannot run plugin %s", plugin)
	}
	return string(out)
}

func initializeConfig(cmd *cobra.Command) error {
	var path string
	dir, err := os.Getwd()
	if err != nil {
		path = ""
	} else {
		path = filepath.Base(dir)
	}

	// load the project config file if it exists, otherwise prompt the user to create one
	return config.LoadConfig(cmd.Flags(), flag.IsSilent(cmd.Flags()), config.ConfigFilename, path)
}
