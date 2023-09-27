package cmd

import (
	"encoding/json"
	"os"

	"github.com/spf13/cobra"
)

var InitRootCmd = &cobra.Command{
	Use:   "init",
	Short: "initialise a project with this plugin",
	Long:  "initialise a project's vision.json file with this plugin's configuration values",
	RunE:  initCmd,
}

var initCmd = func(cmd *cobra.Command, args []string) error {
	const CONFIG = "config" // vision can only accept "config" as the config name
	json.NewEncoder(os.Stdout).Encode(map[string]any{
		CONFIG: map[string]string{
			"key1": "value",
			"key2": "value",
		},
	})
	return nil
}
