package config

const (
	viperServicesPrefix              = "services"
	viperServicesGraphqlKey          = viperServicesPrefix + ".graphql"
	viperServicesGatewayKey          = viperServicesPrefix + ".gateway"
	viperDefaultServicesNamespaceKey = viperServicesPrefix + ".default-namespace"
	viperServicesDirectoryKey        = viperServicesPrefix + ".directory"
)

func GraphqlName() string {
	return v.GetString(viperServicesGraphqlKey)
}

func SetGraphqlName(s string) {
	setAndSave(viperServicesGraphqlKey, s)
}

func GatewayName() string {
	return v.GetString(viperServicesGatewayKey)
}

func SetGatewayName(s string) {
	setAndSave(viperServicesGatewayKey, s)
}

func DefaultNamespace() string {
	return v.GetString(viperDefaultServicesNamespaceKey)
}

func SetDefaultNamespace(s string) {
	setAndSave(viperDefaultServicesNamespaceKey, s)
}

func ServicesDirectory() string {
	return v.GetString(viperServicesDirectoryKey)
}

func SetServicesDirectory(s string) {
	setAndSave(viperServicesDirectoryKey, s)
}
