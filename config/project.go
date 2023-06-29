package config

import (
	"github.com/spf13/viper"
)

const (
	viperProjectPrefix             = "project"
	viperProjectNameKey            = viperProjectPrefix + ".name"
	viperProjectTemplateVersionKey = viperProjectPrefix + ".template-version"
	viperProjectProviderKey        = viperProjectPrefix + ".provider"
	viperProjectDeploymentKey      = viperProjectPrefix + ".deployment"
	viperProjectRemoteKey          = viperProjectPrefix + ".remote"
	viperProjectRegistryKey        = viperProjectPrefix + ".registry"
	viperProjectBranchKey          = viperProjectPrefix + ".branch"
	viperProjectApiVersionKey      = viperProjectPrefix + ".api-version"
)

func ProjectName() string {
	return viper.GetString(viperProjectNameKey)
}

func SetProjectName(s string) {
	setAndSave(viperProjectNameKey, s)
}

func TemplateVersion() string {
	return viper.GetString(viperProjectTemplateVersionKey)
}

func SetTemplateVersion(s string) {
	setAndSave(viperProjectTemplateVersionKey, s)
}

func Provider() string {
	return viper.GetString(viperProjectProviderKey)
}

func SetProvider(s string) {
	setAndSave(viperProjectProviderKey, s)
}

func Deployment() string {
	return viper.GetString(viperProjectDeploymentKey)
}

func SetDeployment(s string) {
	setAndSave(viperProjectDeploymentKey, s)
}

func Remote() string {
	return viper.GetString(viperProjectRemoteKey)
}

func SetRemote(s string) {
	setAndSave(viperProjectRemoteKey, s)
}

func Registry() string {
	return viper.GetString(viperProjectRegistryKey)
}

func SetRegistry(s string) {
	setAndSave(viperProjectRegistryKey, s)
}

func Branch() string {
	return viper.GetString(viperProjectBranchKey)
}

func SetBranch(s string) {
	setAndSave(viperProjectBranchKey, s)
}

func ApiVersion() string {
	return viper.GetString(viperProjectApiVersionKey)
}

func SetApiVersion(s string) {
	setAndSave(viperProjectApiVersionKey, s)
}

func IsDeploymentStandaloneGraphql() bool {
	return Deployment() == DeployStandaloneGraphql
}

func IsDeploymentStandaloneGateway() bool {
	return Deployment() == DeployStandaloneGateway
}

func IsDeploymentPlatform() bool {
	return Deployment() == DeployPlatform
}
