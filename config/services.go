package config

import "github.com/spf13/viper"

const (
	viperServicesPrefix              = "services"
	viperServicesGraphqlKey          = viperServicesPrefix + ".graphql"
	viperServicesGatewayKey          = viperServicesPrefix + ".gateway"
	viperDefaultServicesNamespaceKey = viperServicesPrefix + ".default-namespace"
	viperServicesDirectoryKey        = viperServicesPrefix + ".directory"
)

func GraphqlName() string {
	return viper.GetString(viperServicesGraphqlKey)
}

func SetGraphqlName(s string) {
	setAndSave(viperServicesGraphqlKey, s)
}

func GatewayName() string {
	return viper.GetString(viperServicesGatewayKey)
}

func SetGatewayName(s string) {
	setAndSave(viperServicesGatewayKey, s)
}

func DefaultNamespace() string {
	return viper.GetString(viperDefaultServicesNamespaceKey)
}

func SetDefaultNamespace(s string) {
	setAndSave(viperDefaultServicesNamespaceKey, s)
}

func ServicesDirectory() string {
	return viper.GetString(viperServicesDirectoryKey)
}

func SetServicesDirectory(s string) {
	setAndSave(viperServicesDirectoryKey, s)
}
