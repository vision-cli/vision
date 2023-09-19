package plugins

import (
	"os"
	"path/filepath"

	"github.com/vision-cli/vision/cli"
	"github.com/vision-cli/vision/common/comms"
	"github.com/vision-cli/vision/common/execute"
	"github.com/vision-cli/vision/common/plugins"
	"github.com/vision-cli/vision/config"
	"github.com/vision-cli/vision/flag"
	"github.com/vision-cli/vision/placeholders"

	"github.com/spf13/cobra"
	api_v1 "github.com/vision-cli/api/v1"
)

var UsageQuery = api_v1.PluginRequest{
	Command: api_v1.CommandUsage,
}

func GetCobraCommand(plugin plugins.Plugin, executor execute.Executor) (*cobra.Command, error) {
	usage, err := comms.Call[api_v1.PluginUsageResponse](plugin, &UsageQuery, executor)
	if err != nil {
		return nil, err
	}
	cc := cobra.Command{
		Use:     usage.Use,
		Short:   usage.Short,
		Long:    usage.Long,
		Example: usage.Example,
		Run: func(cmd *cobra.Command, args []string) {
			err := initializeConfig(cmd, usage.RequiresConfig)
			if err != nil && usage.RequiresConfig {
				cli.Fatalf("cannot initialize config: %v", err)
			}
			p, err := placeholders.NewPlaceholders(cmd.Flags(), ".", "default", "", "")
			if err != nil {
				cli.Fatalf("cannot initialize placeholders: %v", err)
			}
			response, err := comms.Call[api_v1.PluginResponse](plugin, &api_v1.PluginRequest{
				Command:      api_v1.CommandRun,
				Args:         args,
				Flags:        []api_v1.PluginFlag{},
				Placeholders: *p,
			}, executor)
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
	return &cc, nil
}

func initializeConfig(cmd *cobra.Command, requireConfig bool) error {
	var path string
	dir, err := os.Getwd()
	if err != nil {
		path = ""
	} else {
		path = filepath.Base(dir)
	}

	// load the project config file if it exists, otherwise prompt the user to create one
	return config.LoadConfig(cmd.Flags(), flag.IsSilent(cmd.Flags()), config.ConfigFilename, path, requireConfig)
}
