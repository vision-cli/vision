package config

const (
	uniqueStr       = "unique-str"
	uniqueStrLen    = 6
	namifyLen       = 20
	ConfigFilename  = "jarvis"
	configExtension = ".json"
)

const (
	defaultTempalateVersion = "v1"
	defaultGatewayName      = "gateway"
	defaultGraphqlName      = "graphql"
	defaultDefaultNamespace = "default"
	defaultServicesDir      = "services"
	defaultInfraDir         = "infra"
	defaultRegistry         = "gcr.io"
	defaultBranch           = "master"
	defaultApiVersion       = "v1"
	defaultDeployment       = DeployStandaloneGraphql
)

const (
	DeployStandaloneGraphql = "standalone-graphql"
	DeployStandaloneGateway = "standalone-gateway"
	DeployPlatform          = "platform"
)

const (
	FlagSilent     = "silent"
	FlagForce      = "force"
	FlagRemote     = "remote"
	FlagBranch     = "branch"
	FlagRegistry   = "registry"
	FlagNamespace  = "namespace"
	FlagApiVersion = "version"
	FlagDeployment = "deployment"
	FlagTemplate   = "template"
)
