package cmd

import (
	"encoding/json"
	"os"

	"github.com/spf13/cobra"
)

var VersionRootCmd = &cobra.Command{
	Use:   "version",
	Short: "the plugin version",
	Long:  "ditto",
	RunE:  versionCmd,
}

var versionCmd = func(cmd *cobra.Command, args []string) error {
	json.NewEncoder(os.Stdout).Encode(map[string]string{
		"sem_ver": "0.0.1",
	})

	return nil
}
