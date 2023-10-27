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
	pd := PluginData{
		PluginConfig: PluginConfig{
			PluginName: "sample-plugin",
			ModuleName: "github.com/my-org/my-plugin",
		},
	}

	err := json.NewEncoder(os.Stdout).Encode(pd)

	if err != nil {
		return fmt.Errorf("encoding json: %w", err)
	}
	return nil
}
