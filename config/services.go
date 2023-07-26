package config

const (
	viperServicesPrefix              = "services"
	viperServicesGraphqlKey          = viperServicesPrefix + ".graphql"
	viperServicesGatewayKey          = viperServicesPrefix + ".gateway"
	viperDefaultServicesNamespaceKey = viperServicesPrefix + ".default-namespace"
	viperServicesDirectoryKey        = viperServicesPrefix + ".directory"
)

func GraphqlName() string {
	return v.GetStringOrDefault(viperServicesGraphqlKey, defaultGraphqlName)
}

func SetGraphqlName(s string) {
	setAndSave(viperServicesGraphqlKey, s)
}

func GatewayName() string {
	return v.GetStringOrDefault(viperServicesGatewayKey, defaultGatewayName)
}

func SetGatewayName(s string) {
	setAndSave(viperServicesGatewayKey, s)
}

func DefaultNamespace() string {
	return v.GetStringOrDefault(viperDefaultServicesNamespaceKey, defaultDefaultNamespace)
}

func SetDefaultNamespace(s string) {
	setAndSave(viperDefaultServicesNamespaceKey, s)
}

func ServicesDirectory() string {
	return v.GetStringOrDefault(viperServicesDirectoryKey, defaultServicesDir)
}

func SetServicesDirectory(s string) {
	setAndSave(viperServicesDirectoryKey, s)
}
