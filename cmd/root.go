package cmd

import (
	_ "embed"
	"encoding/json"
	"io"
	"os"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/vision-cli/vision/cmd/doctor"
	initialise "github.com/vision-cli/vision/cmd/init"
	"github.com/vision-cli/vision/cmd/list"
	"github.com/vision-cli/vision/internal/plugin"
)

// Finds available plugins and initialises them into commands
func init() {
	rootCmd.AddCommand(initialise.RootCmd)
	rootCmd.AddCommand(doctor.RootCmd)
	rootCmd.AddCommand(list.RootCmd)
	rootCmd.Flags().AddFlagSet(initVisionFlags())
	plugins := plugin.Find()
	for _, plugin := range plugins {
		cmd, err := createCommand(plugin)
		if err != nil {
			// TODO(steve): handle broken commands, maybe a vision doctor command???
			log.Error(plugin.FullPath, err)
			continue
		}
		rootCmd.AddCommand(cmd)
	}
}

// createCommand takes in a plugin and returns a cobra command to interact with that plugin
func createCommand(p plugin.Plugin) (*cobra.Command, error) {
	exe := plugin.NewExecutor(p.FullPath)
	info, err := exe.Info()
	if err != nil {
		return nil, err
	}
	version, err := exe.Version()
	if err != nil {
		return nil, err
	}

	cobraCmd := &cobra.Command{
		Use:     p.Name,
		Version: version.SemVer,
		Short:   info.ShortDescription,
		Long:    info.LongDescription,
		RunE:    createPluginCommandHandler(p),
	}

	return cobraCmd, nil
}

func createPluginCommandHandler(p plugin.Plugin) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 { // prevents index out of range
			log.Warnf("No argument provided. Try: \n\t\n vision %v -v", cmd.Use)
		}
		exe := plugin.NewExecutor(p.FullPath)
		switch args[0] {
		case "init":
			i, err := exe.Init()
			if err != nil {
				return err
			}
			// TODO merge into vision config
			mergeConfigs(p.Name, i.Config)
		}
		return nil
	}
}

func mergeConfigs(pluginName string, config any) error {
	writeSuccess := false
	isTruncated := false
	f, err := os.OpenFile("vision.json", os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	defer func() {
		// defensive coding to keep a clone of the original data
		var originalData []byte
		copy(originalData, b)
		if !writeSuccess && isTruncated {
			_, err = f.Write(originalData)
			if err != nil {
				log.Fatal(err)
			}
		}
	}()

	var configBytes map[string]any
	err = json.Unmarshal(b, &configBytes)
	if err != nil {
		return err
	}

	configBytes[pluginName] = config
	err = f.Truncate(0)
	if err != nil {
		return err
	}
	isTruncated = true
	_, err = f.Seek(0, 0)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	err = enc.Encode(configBytes)
	if err != nil {
		return err
	}

	writeSuccess = true
	return nil
}

func initVisionFlags() *pflag.FlagSet {
	fs := pflag.NewFlagSet("vision", 1)
	return fs
}

//go:embed example.txt
var exampleText string

var rootCmd = &cobra.Command{
	Use:     "vision",
	Short:   "A developer productivity tool",
	Long:    `Vision is a tool to create microservice platforms and microservice scaffolding code`,
	Example: exampleText,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Error(err)
	}
}
