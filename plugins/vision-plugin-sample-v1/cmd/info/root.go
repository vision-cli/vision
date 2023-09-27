package cmd

import (
	_ "embed"
	"encoding/json"
	"os"

	"github.com/spf13/cobra"
)

var InfoRootCmd = &cobra.Command{
	Use:   "info",
	Short: "return info about the plugin",
	Long:  "ditto",
	RunE:  sampleCmd,
}

//go:embed info.txt
var infoOutput string

var sampleCmd = func(cmd *cobra.Command, args []string) error {
	json.NewEncoder(os.Stdout).Encode(map[string]string{
		"short_description": infoOutput,
		"long_description":  infoOutput,
	})

	return nil
}
