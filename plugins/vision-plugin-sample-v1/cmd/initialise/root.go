package cmd

import (
	"encoding/json"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func init() {
	InitRootCmd.Flags().AddFlagSet(initialiseFlags())
}

var (
	configValue string
)

func initialiseFlags() *pflag.FlagSet {
	fs := pflag.NewFlagSet("init", 1)
	fs.StringVarP(&configValue, "confValue", "c", "", "set the config value")
	return fs
}

var InitRootCmd = &cobra.Command{
	Use:   "init",
	Short: "initialise a project with this plugin",
	Long:  "initialise a project's vision.json file with this plugin's configuration values",
	RunE:  initCmd,
}

var initCmd = func(cmd *cobra.Command, args []string) error {
	const CONFIG = "config" // vision only accepts "config" as the config name
	json.NewEncoder(os.Stdout).Encode(map[string]any{
		CONFIG: map[string]string{
			"key1": configValue,
			"key2": "value",
		},
	})
	return nil
}
