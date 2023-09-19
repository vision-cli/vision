package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func init() {
	RootCmd.Flags().AddFlagSet(initFlags())
}

// left as an example of flags for future reference
// var (
// 	MyBool bool
// )

func initFlags() *pflag.FlagSet {
	fs := pflag.NewFlagSet("init", 1)
	// placefholder for future flags
	// fs.BoolVarP(&MyBool, "tf", "t", true, "test flag")
	return fs
}

var RootCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialise a new vision project",
	Long:  "Create a new vision project and initialise default config values for vision and installed plugins",
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, err := os.Getwd()
		if err != nil {
			return err
		}
		// take the current dir name to use as project name
		projectName := filepath.Base(dir)
		// check if file exists, if so do nothing leaving info log for user
		configFilePath := filepath.Join(dir, configFileName)
		if _, err := os.Stat(configFilePath); errors.Is(err, os.ErrNotExist) {
			// file does not exist so we create it
			err = createDefaultConfig(projectName)
			if err != nil {
				return err
			}
			log.Info("successfully initialised vision")
			return nil
		}
		// file already exists, let user know
		log.Info("config file already exists")
		return nil
	},
}

const configFileName = "vision.json"

type VisionConfig struct {
	ProjectName string `json:"project_name"`
}

// create a default json file with basic info as defined in the config model.
// TODO(steve): generate default config for each installed plugin
func createDefaultConfig(projectName string) error {
	f, err := os.Create(configFileName)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(VisionConfig{
		ProjectName: projectName,
	})
}
