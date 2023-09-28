package cmd

import (
	"encoding/json"
	"os"

	"github.com/spf13/cobra"
	// "github.com/spf13/pflag"
)

func init() {
	InitRootCmd.PersistentFlags().StringVarP(&configValue, "confValue", "c", "", "set the config value")
	InitRootCmd.PersistentFlags().StringVarP(&configValue, "testValue", "t", "", "set the test value")
	InitRootCmd.PersistentFlags().StringVarP(&configValue, "anotherValue", "a", "", "set another value")
	InitRootCmd.PersistentFlags().BoolVarP(&boolValue, "boolValue", "b", false, "set a bool value")
}

var (
	configValue string
	boolValue   bool
)

// func initialiseFlags() *pflag.FlagSet {
// 	fs := pflag.NewFlagSet("init", 1)
// 	fs.
// 	return fs
// }

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
		},
	})

	return nil
}
