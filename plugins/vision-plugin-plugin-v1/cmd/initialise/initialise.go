package initialise

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type PluginConfig struct {
	PluginName string `json:"plugin_name"`
	ModuleName string `json:"module_name"`
	Command    string `json:"command"`
}

type PluginData struct {
	PluginConfig PluginConfig `json:"plugin"`
}

type visionPluginData struct {
	PluginConfig PluginConfig `json:"config"`
}

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "initialise a project with this plugin",
	Long:  "initialise a project's vision.json file with this plugin's configuration values",
	RunE:  runCommand,
}

func runCommand(cmd *cobra.Command, args []string) error {

	pd := visionPluginData{
		PluginConfig: PluginConfig{
			PluginName: "sample-plugin",
			ModuleName: "github.com/my-org/my-plugin",
			Command:    "changeme",
		},
	}

	err := json.NewEncoder(os.Stdout).Encode(pd)

	if err != nil {
		return fmt.Errorf("encoding json: %w", err)
	}
	return nil
}
