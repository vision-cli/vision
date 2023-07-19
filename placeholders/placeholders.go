package placeholders

import (
	"fmt"
	"net/url"
	"path/filepath"

	"github.com/spf13/pflag"

	api_v1 "github.com/vision-cli/api/v1"
	"github.com/vision-cli/common/cases"
	"github.com/vision-cli/vision/config"
)

func getValueFromFlagSetOrConfig(flagSet *pflag.FlagSet, flagName string, configFunc func() string) string {
	if val, err := flagSet.GetString(flagName); err == nil && val != "" {
		return val
	}
	return configFunc()
}

func NewProjectPlaceholders(flagSet *pflag.FlagSet,
	projectRoot, projectName string) (*api_v1.PluginPlaceholders, error) {
	return NewPlaceholders(flagSet, projectRoot, projectName, "", "")
}

func NewServicePlaceholders(flagSet *pflag.FlagSet,
	projectRoot, moduleName, serviceName string) (*api_v1.PluginPlaceholders, error) {
	return NewPlaceholders(flagSet, projectRoot, config.ProjectName(), moduleName, serviceName)
}

func NewDefaultServicePlaceholders(flagSet *pflag.FlagSet,
	projectRoot, serviceName string) (*api_v1.PluginPlaceholders, error) {
	return NewPlaceholders(flagSet, projectRoot, config.ProjectName(), "", serviceName)
}

func NewPlaceholders(flagSet *pflag.FlagSet,
	rawProjectRoot,
	rawProjectName,
	rawServiceNamespace,
	rawServiceName string) (*api_v1.PluginPlaceholders, error) {
	// Project name is snake case for use as a variable
	projectName := cases.Kebab(rawProjectName)
	// replace the raw project name if there is a config project name
	if config.ProjectName() != "" {
		projectName = cases.Kebab(config.ProjectName())
	}

	// Project directory is kebab case for use as a folder name
	projectDirectory := cases.Kebab(rawProjectName)

	// Service name is snake case for use as a variable
	serviceName := cases.Snake(rawServiceName)

	// Service namespace is (in priority order): passed parameter, flag value, default value
	serviceNamespace := rawServiceNamespace
	if serviceNamespace == "" {
		if flagSet.Changed(config.FlagNamespace) {
			serviceNamespace = flagSet.Lookup(config.FlagNamespace).Value.String()
		} else {
			serviceNamespace = config.DefaultNamespace()
		}
	}

	registry := getValueFromFlagSetOrConfig(flagSet, config.FlagRegistry, config.Registry)
	remote := getValueFromFlagSetOrConfig(flagSet, config.FlagRemote, config.Remote)
	branch := getValueFromFlagSetOrConfig(flagSet, config.FlagBranch, config.Branch)
	version := getValueFromFlagSetOrConfig(flagSet, config.FlagApiVersion, config.ApiVersion)
	servicesDir := config.ServicesDirectory()

	if remote == "" {
		return nil, fmt.Errorf("remote cannot be empty, please provide a remote with the -r --remote flag or set it in the config file")
	}

	projectFqn, err := url.JoinPath(remote, projectName)
	if err != nil {
		return nil, err
	}
	servicesFqn, err := url.JoinPath(projectFqn, config.ServicesDirectory(), serviceNamespace)
	if err != nil {
		return nil, err
	}
	libsfqn, err := url.JoinPath(projectFqn, "libs")
	if err != nil {
		return nil, err
	}
	servicefqn, err := url.JoinPath(servicesFqn, serviceName)
	if err != nil {
		return nil, err
	}

	graphqlServiceName := cases.Snake(config.GraphqlName())
	gatewayServiceName := cases.Snake(config.GatewayName())
	gatewayfqn, err := url.JoinPath(servicesFqn, gatewayServiceName)
	if err != nil {
		return nil, err
	}
	graphqlfqn, err := url.JoinPath(servicesFqn, graphqlServiceName)
	if err != nil {
		return nil, err
	}

	deployment := config.Deployment()

	return &api_v1.PluginPlaceholders{
		// project
		ProjectRoot:      rawProjectRoot,
		ProjectName:      projectName,
		ProjectDirectory: projectDirectory,
		ProjectFqn:       projectFqn,
		Registry:         registry,
		Remote:           remote,
		Branch:           branch,
		Version:          version,
		// services
		ServicesFqn:        servicesFqn,
		ServicesDirectory:  servicesDir,
		GatewayServiceName: gatewayServiceName,
		GatewayFqn:         gatewayfqn,
		GraphqlServiceName: graphqlServiceName,
		GraphqlFqn:         graphqlfqn,
		LibsFqn:            libsfqn,
		LibsDirectory:      filepath.Join(rawProjectRoot, "libs"),
		// service
		ServiceNamespace:          serviceNamespace,
		ServiceVersionedNamespace: serviceNamespace + "." + version,
		ServiceName:               serviceName,
		ServiceFqn:                servicefqn,
		ServiceDirectory:          filepath.Join(projectDirectory, servicesDir, serviceNamespace, rawServiceName),
		// infra
		InfraDirectory: config.InfraDirectory(),
		// messaging
		ProtoPackage: serviceNamespace + "." + rawServiceName + "." + version,
		// graphql
		Deployment: deployment,
	}, nil
}
