//go:build vision_dev

package plugins

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/vision-cli/vision/common/marshal"
	"github.com/vision-cli/vision/common/plugins"
	"github.com/vision-cli/vision/common/tmpl"
	pluginPlugin "github.com/vision-cli/vision/core/plugin"
	projectPlugin "github.com/vision-cli/vision/core/project"
	"github.com/vision-cli/vision/execute"

	api_v1 "github.com/vision-cli/api/v1"
	gatewayPlugin "github.com/vision-cli/vision-plugin-gateway-v1/plugin"
	graphqlPlugin "github.com/vision-cli/vision-plugin-graphql-v1/plugin"
	infraPlugin "github.com/vision-cli/vision-plugin-infra-v1/plugin"
	servicePlugin "github.com/vision-cli/vision-plugin-service-v1/plugin"
)

func init() {
	plugins.InternalPlugins = append(plugins.InternalPlugins,
		plugins.Plugin{
			Name:            "vision-plugin-gateway-v1",
			PluginPath:      "__vision-internal-plugin",
			InternalCommand: gatewayPlugin.Handle,
		},
		plugins.Plugin{
			Name:            "vision-plugin-graphql-v1",
			PluginPath:      "__vision-internal-plugin",
			InternalCommand: graphqlPlugin.Handle,
		},
		plugins.Plugin{
			Name:            "vision-plugin-infra-v1",
			PluginPath:      "__vision-internal-plugin",
			InternalCommand: infraPlugin.Handle,
		},
		plugins.Plugin{
			Name:            "vision-plugin-plugin-v1",
			PluginPath:      "__vision-internal-plugin",
			InternalCommand: pluginPlugin.Handle,
		},
		plugins.Plugin{
			Name:            "vision-plugin-project-v1",
			PluginPath:      "__vision-internal-plugin",
			InternalCommand: projectPlugin.Handle,
		},
		plugins.Plugin{
			Name:            "vision-plugin-service-v1",
			PluginPath:      "__vision-internal-plugin",
			InternalCommand: servicePlugin.Handle,
		},
		plugins.Plugin{
			Name:            "vision-plugin-debug-v1",
			PluginPath:      "__vision-internal-plugin",
			InternalCommand: printPluginRequest,
		},
	)
}

func printPluginRequest(input string, _ execute.Executor, _ tmpl.TmplWriter) string {
	req, err := marshal.Unmarshal[api_v1.PluginRequest](input)
	if err != nil {
		return errorResponse(err)
	}
	result := ""
	switch req.Command {
	case api_v1.CommandUsage:
		result, err = marshal.Marshal[api_v1.PluginUsageResponse](debugPluginUsage)
		if err != nil {
			return errorResponse(err)
		}
	case api_v1.CommandConfig:
		result, err = marshal.Marshal[api_v1.PluginConfigResponse](debugPluginConfig)
		if err != nil {
			return errorResponse(err)
		}
	case api_v1.CommandRun:
		formattedRequest, err := indent(input)
		fmt.Println("Raw plugin request:", formattedRequest)
		resp := api_v1.PluginResponse{
			Result: "SUCCESS!",
			Error:  "",
		}
		result, err = marshal.Marshal[api_v1.PluginResponse](resp)
		if err != nil {
			return errorResponse(err)
		}
	default:
		return errorResponse(errors.New("unknown command"))
	}
	return result
}

var debugPluginUsage = api_v1.PluginUsageResponse{
	Version:        "0.0.0",
	Use:            "printRequest",
	Short:          "print plugin request",
	Long:           "print the request sent to a plugin",
	Example:        "vision printRequest",
	Subcommands:    []string{"debug"},
	Flags:          []api_v1.PluginFlag{},
	RequiresConfig: false,
}

var debugPluginConfig = api_v1.PluginConfigResponse{
	Defaults: []api_v1.PluginConfigItem{},
}

func errorResponse(err error) string {
	res, err := marshal.Marshal[api_v1.PluginResponse](api_v1.PluginResponse{
		Result: "",
		Error:  err.Error(),
	})
	if err != nil {
		panic(err.Error())
	}
	return res
}

func indent(rawJson string) (string, error) {
	var buf bytes.Buffer
	err := json.Indent(&buf, []byte(rawJson), "", "    ")
	if err != nil {
		return "", err
	}
	return string(buf.Bytes()), nil
}
