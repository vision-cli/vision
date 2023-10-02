package initialise

import (
	"encoding/json"
	"os"

	"github.com/spf13/cobra"
)

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "initialise a project with this plugin",
	Long:  "initialise a project's vision.json file with this plugin's configuration values",
	Run: func(cmd *cobra.Command, args []string) {
		json.NewEncoder(os.Stdout).Encode(map[string]any{
			"config": map[string]any{
				"name": "sample-plugin",
			},
		})
	}}
