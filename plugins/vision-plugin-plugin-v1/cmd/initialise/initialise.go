package initialise

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

type pluginConfig struct {
	PluginName string `json:"plugin_name"`
	ModuleName string `json:"module_name"`
}

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "initialise a project with this plugin",
	Long:  "initialise a project's vision.json file with this plugin's configuration values",
	RunE:  runCommand,
}

func runCommand(cmd *cobra.Command, args []string) error {
	path := args[0]
	err := os.MkdirAll(path, 0644)
	if err != nil {
		return fmt.Errorf("creating init dir: %w", err)
	}

	f, err := os.Create(filepath.Join(path, "vision.json"))
	if err != nil {
		return fmt.Errorf("creating vision.json: %w", err)
	}

	pc := pluginConfig{
		PluginName: "sample-plugin",
		ModuleName: "github.com/my-org/my-plugin",
	}

	return json.NewEncoder(f).Encode(pc)
}
