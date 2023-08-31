package config

const (
	uniqueStr       = "unique-str"
	uniqueStrLen    = 6
	namifyLen       = 20
	ConfigFilename  = "vision"
	configExtension = ".json"
)

const (
	defaultTemplateVersion  = "v1"
	defaultGatewayName      = "gateway"
	defaultGraphqlName      = "graphql"
	defaultDefaultNamespace = "default"
	defaultServicesDir      = "services"
	defaultInfraDir         = "infra"
	defaultBranch           = "master"
	defaultApiVersion       = "v1"
	defaultDeployment       = DeployStandaloneGraphql
	defaultRegistry         = "gcr.io"
)

const (
	DeployStandaloneGraphql = "standalone-graphql"
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

type DefaultConfig struct {
	prompt   string
	def      string
	flagname string
	setter   func(string)
}

var defaultConfigsWithFlags = []DefaultConfig{
	{"Default remote (e.g. github.com/<company-name>):", "", FlagRemote, SetRemote},
	{"Default deployment:", defaultDeployment, FlagDeployment, SetDeployment},
	{"Default service namespace:", defaultDefaultNamespace, FlagNamespace, SetDefaultNamespace},
	{"Default registry (e.g. gcr.io):", "", FlagRegistry, SetRegistry},
	{"Default branch:", defaultBranch, FlagBranch, SetBranch},
	{"Default api version:", defaultApiVersion, FlagApiVersion, SetApiVersion},
}

var defaultConfigs = []DefaultConfig{
	{"Template version:", defaultTemplateVersion, "", SetTemplateVersion},
	{"Default gateway service name:", defaultGatewayName, "", SetGatewayName},
	{"Default graphql service name:", defaultGraphqlName, "", SetGraphqlName},
	{"Default services directory:", defaultServicesDir, "", SetServicesDirectory},
	{"Default infra directory:", defaultInfraDir, "", SetInfraDirectory},
}
