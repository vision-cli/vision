package cmd

import (
	"encoding/json"
	"os"

	"github.com/spf13/cobra"
	// "github.com/spf13/pflag"
)

func init() {
	InitCmd.PersistentFlags().StringVarP(&configValue, "confValue", "c", "", "set the config value")
	InitCmd.PersistentFlags().StringVarP(&configValue, "testValue", "t", "", "set the test value")
	InitCmd.PersistentFlags().StringVarP(&configValue, "anotherValue", "a", "", "set another value")
	InitCmd.PersistentFlags().BoolVarP(&boolValue, "boolValue", "b", false, "set a bool value")
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

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "initialise a project with this plugin",
	Long:  "initialise a project's vision.json file with this plugin's configuration values",
	Run: func(cmd *cobra.Command, args []string) {
		json.NewEncoder(os.Stdout).Encode(map[string]any{
			"config": map[string]string{ // vision only accepts "config" as the config name
				"key1": configValue,
			},
		})
	}}
