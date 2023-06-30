package plugins

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	api_v1 "github.com/vision-cli/api/v1"
	"github.com/vision-cli/vision/cli"
	"github.com/vision-cli/vision/config"
	"github.com/vision-cli/vision/execute"
	"github.com/vision-cli/vision/flag"
	"github.com/vision-cli/vision/placeholders"
)

func GetCobraCommand(plugin string) *cobra.Command {
	pluginPath := filepath.Join(goBinPath(), plugin)
	osExecutor := execute.NewOsExecutor()
	usage, err := Call[api_v1.PluginUsageResponse](pluginPath, &UsageQuery, osExecutor)
	if err != nil {
		cli.Fatalf(err.Error())
	}
	cc := cobra.Command{
		Use:     usage.Use,
		Short:   usage.Short,
		Long:    usage.Long,
		Example: usage.Example,
		Run: func(cmd *cobra.Command, args []string) {
			if usage.RequiresConfig {
				if err := initializeConfig(cmd); err != nil {
					cli.Fatalf("cannot initialize config: %v", err)
				}
			}
			placeholders := placeholders.NewPlaceholders(cmd.Flags(), ".", "default", "", "")
			response, err := Call[api_v1.PluginResponse](pluginPath, &api_v1.PluginRequest{
				Command:      api_v1.CommandRun,
				Args:         args,
				Flags:        []api_v1.PluginFlag{},
				Placeholders: placeholders,
			}, osExecutor)
			if err != nil {
				cli.Fatalf(err.Error())
			}
			if response.Error != "" {
				cli.Fatalf(response.Error)
			}
			cli.Infof(response.Result)
		},
	}
	cc.Flags().AddFlagSet(flag.ConfigFlagset())
	return &cc
}

func initializeConfig(cmd *cobra.Command) error {
	var path string
	dir, err := os.Getwd()
	if err != nil {
		path = ""
	} else {
		path = filepath.Base(dir)
	}

	// load the project config file if it exists, otherwise prompt the user to create one
	return config.LoadConfig(cmd.Flags(), flag.IsSilent(cmd.Flags()), config.ConfigFilename, path)
}
