package initialise

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

type PluginConfig struct {
	PluginName string `json:"plugin_name"`
	ModuleName string `json:"module_name"`
	GoVersion  string `json:"go_version"`
}

type PluginData struct {
	PluginConfig PluginConfig `json:"config"`
}

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "initialise a project with this plugin",
	Long:  "initialise a project's vision.json file with this plugin's configuration values",
	RunE:  runCommand,
}

func runCommand(cmd *cobra.Command, args []string) error {
	gv, err := getLatestGoVersion()
	if err != nil {
		return fmt.Errorf("getting go version: %w", err)
	}

	pd := PluginData{
		PluginConfig: PluginConfig{
			PluginName: "sample-plugin",
			ModuleName: "github.com/my-org/my-plugin",
			GoVersion:  gv,
		},
	}

	err = json.NewEncoder(os.Stdout).Encode(pd)

	if err != nil {
		return fmt.Errorf("encoding json: %w", err)
	}
	return nil
}

func getLatestGoVersion() (string, error) {
	cmd := "curl 'https://go.dev/VERSION?m=text'"
	b, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return "", fmt.Errorf("curling Go version: %w", err)
	}

	goVersion := string(b)[2:8]

	return goVersion, nil
}
