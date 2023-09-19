package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/charmbracelet/log"
	"github.com/vision-cli/vision/config"
	"github.com/vision-cli/vision/flag"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new config file",
	Long:  `Run the config creation wizard to create a new config file`,
	Args:  cobra.MaximumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		var path string
		dir, err := os.Getwd()
		if err != nil {
			path = ""
		} else {
			path = filepath.Base(dir)
		}

		if err := config.LoadConfig(cmd.Flags(), flag.IsSilent(cmd.Flags()), config.ConfigFilename, path, true); err != nil {
			log.Error("cannot create config file: %v", err)
		}
	},
}
