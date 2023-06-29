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
)

var UsageQuery = api_v1.PluginRequest{
	Command:      api_v1.CommandUsage,
	Args:         "",
	Placeholders: "",
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
	usageQuery, err := MarshalRequest(UsageQuery)
	if err != nil {
		cli.Fatalf("Cannot marshal usage query for plugin %s", plugin)
	}
	response := callPlugin(pluginPath, usageQuery)
	usage, err := UnmarshalResponse[api_v1.PluginUsageResponse](response)
	if err != nil {
		cli.Fatalf("Cannot marshal usage query for plugin %s", plugin)
	}
	return &cobra.Command{
		Use:     usage.Use,
		Short:   usage.Short,
		Long:    usage.Long,
		Example: usage.Example,
		Run: func(cmd *cobra.Command, args []string) {
			runQuery, err := MarshalRequest(api_v1.PluginRequest{
				Command:      api_v1.CommandRun,
				Args:         "",
				Placeholders: "",
			})
			if err != nil {
				cli.Fatalf("Cannot marshal run query for plugin %s", plugin)
			}
			response := callPlugin(pluginPath, runQuery)
			result, err := UnmarshalResponse[api_v1.PluginRunResponse](response)
			if err != nil {
				cli.Fatalf("Cannot unmarshal response from plugin %s", plugin)
			}
			cli.Infof(result.Result)
		},
	}
}

func UnmarshalResponse[T any](reqStr string) (T, error) {
	var req T
	err := json.Unmarshal([]byte(reqStr), &req)
	if err != nil {
		return req, err
	}
	return req, nil
}

func MarshalRequest[T any](resp T) (string, error) {
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
	// eventialluy remove 'example' as well
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
