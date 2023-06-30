package config

const (
	uniqueStr       = "unique-str"
	uniqueStrLen    = 6
	namifyLen       = 20
	ConfigFilename  = "vision"
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

type DefaultConfig struct {
	prompt   string
	def      string
	flagname string
	setter   func(string)
}

var defaultConfigsWithFlags = []DefaultConfig{
	{"Default remote (e.g. github.com/<company-name>/<project>):", "", FlagRemote, SetRemote},
	{"Default deployment:", defaultDeployment, FlagDeployment, SetDeployment},
	{"Default service namespace:", defaultDefaultNamespace, FlagNamespace, SetDefaultNamespace},
	{"Default registry:", defaultRegistry, FlagRegistry, SetRegistry},
	{"Default branch:", defaultBranch, FlagBranch, SetBranch},
	{"Default api version:", defaultApiVersion, FlagApiVersion, SetApiVersion},
}

var defaultConfigs = []DefaultConfig{
	{"Template version:", defaultTempalateVersion, "", SetTemplateVersion},
	{"Default gateway service name:", defaultGatewayName, "", SetGatewayName},
	{"Default graphql service name:", defaultGraphqlName, "", SetGraphqlName},
	{"Default services directory:", defaultServicesDir, "", SetServicesDirectory},
	{"Default infra directory:", defaultInfraDir, "", SetInfraDirectory},
}
