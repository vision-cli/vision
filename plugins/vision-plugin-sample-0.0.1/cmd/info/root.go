package cmd

import (
	_ "embed"
	"encoding/json"
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "info",
	Short: "the plugin that is used to generate new plugins",
	Long:  "the plugin that is used to generate new plugins",
	RunE:  infoCmd,
}

//go:embed info.txt
var infoOut string

var infoCmd = func(cmd *cobra.Command, args []string) error {

	json.NewEncoder(os.Stdout).Encode(infoOut)
	return nil
}
