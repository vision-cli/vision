package config

import (
	"encoding/json"
	"errors"
	"fmt"
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
var (
	projectName string
)

func initFlags() *pflag.FlagSet {
	fs := pflag.NewFlagSet("init", 1)
	fs.StringVarP(&projectName, "project", "p", "", "project name")
	return fs
}

var RootCmd = &cobra.Command{
	Use:   "vision init [PROJECT_NAME] [OPTIONS]",
	Short: "Initialise a new vision project",
	Long:  "Create a new vision project and initialise default config values for vision and installed plugins",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			log.Info("usage: vision init [projectname]")
			return errors.New("unexpected arguments")
		}

		// take the current dir name to use as project name
		dir, err := os.Getwd()
		if err != nil {
			return err
		}

		var projectDir string
		if len(args) == 0 {
			projectDir = ""
			if projectName == "" { // flag not set by user so we set project name here
				projectName = filepath.Base(dir)
			}
		} else {
			projectDir = args[0]
			if projectName == "" { // flag not set by user so we set project name here
				projectName = args[0]
			}
		}

		configFilePath := filepath.Join(dir, projectDir, configFileName)
		err = createDefaultConfig(configFilePath, projectName, projectDir)
		if err != nil {
			return err
		}
		log.Info("successfully initialised vision")
		return nil
	},
}

const configFileName = "vision.json"

type VisionConfig struct {
	ProjectName string `json:"project_name"`
}

// create a default json file with basic info as defined in the config model.
// if the projectDir is not an empty string, create the directory as well as the file
// TODO(steve): generate default config for each installed plugin
func createDefaultConfig(configFilePath, projectName, projectDir string) error {
	if projectDir != "" {
		if err := os.MkdirAll(filepath.Dir(configFilePath), os.ModePerm); err != nil {
			return fmt.Errorf("error creating directory: %v", err)
		}
	}

	// check if file exists, create if not
	f, err := os.OpenFile(configFilePath, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if errors.Is(err, os.ErrExist) {
		log.Info("config file already exists")
		return nil
	}
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
