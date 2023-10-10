package version

import (
	"encoding/json"
	"os"

	"github.com/spf13/cobra"
)

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "the plugin version",
	Long:  "ditto",
	Run: func(cmd *cobra.Command, args []string) {
		json.NewEncoder(os.Stdout).Encode(map[string]string{
			"sem_ver": "v0.0.1",
		})
	},
}
