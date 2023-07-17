//go:build vision_dev

package plugins

import (
	"github.com/vision-cli/common/plugins"

	gatewayPlugin "github.com/vision-cli/vision-plugin-gateway-v1/plugin"
	graphqlPlugin "github.com/vision-cli/vision-plugin-graphql-v1/plugin"
	infraPlugin "github.com/vision-cli/vision-plugin-infra-v1/plugin"
	pluginPlugin "github.com/vision-cli/vision-plugin-plugin-v1/plugin"
	projectPlugin "github.com/vision-cli/vision-plugin-project-v1/plugin"
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
	)
}
