package initialise

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

type PluginConfig struct {
	PluginName string `json:"plugin_name"`
	ModuleName string `json:"module_name"`
}

type PluginData struct {
	PluginData PluginConfig `json:"plugin_plugin"`
}

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "initialise a project with this plugin",
	Long:  "initialise a project's vision.json file with this plugin's configuration values",
	RunE:  runCommand,
}

type visionJson struct {
	Contents map[any]any
}

func runCommand(cmd *cobra.Command, args []string) error {
	var path string
	if args[0] == "" {
		path = "."
	} else {
		path = args[0]
	}
	err := os.MkdirAll(strings.TrimSuffix(path, "vision.json"), os.ModePerm)
	if err != nil {
		return fmt.Errorf("creating init dir: %w", err)
	}

	// TODO (luke): enable PluginName and ModuleName to be read from config
	pd := PluginData{
		PluginData: PluginConfig{
			PluginName: "sample-plugin",
			ModuleName: "github.com/my-org/my-plugin",
		},
	}

	// vPath := filepath.Join(path, "vision.json")
	// f, err := os.Open(vPath)
	// if err != nil {
	// 	return fmt.Errorf("creating vision.json: %w", err)
	// }
	// defer f.Close()

	// b, err := io.ReadAll(f)
	// if err != nil {
	// 	return fmt.Errorf("reading vision.json: %w", err)
	// }

	// err = json.Unmarshal(b)

	err = json.NewEncoder(os.Stdout).Encode(pd)

	if err != nil {
		return fmt.Errorf("encoding json: %w", err)
	}
	return nil
}
