package config

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
	return v.GetString(viperProjectNameKey)
}

func SetProjectName(s string) {
	setAndSave(viperProjectNameKey, s)
}

func TemplateVersion() string {
	return v.GetString(viperProjectTemplateVersionKey)
}

func SetTemplateVersion(s string) {
	setAndSave(viperProjectTemplateVersionKey, s)
}

func Provider() string {
	return v.GetString(viperProjectProviderKey)
}

func SetProvider(s string) {
	setAndSave(viperProjectProviderKey, s)
}

func Deployment() string {
	return v.GetString(viperProjectDeploymentKey)
}

func SetDeployment(s string) {
	setAndSave(viperProjectDeploymentKey, s)
}

func Remote() string {
	return v.GetString(viperProjectRemoteKey)
}

func SetRemote(s string) {
	setAndSave(viperProjectRemoteKey, s)
}

func Registry() string {
	return v.GetString(viperProjectRegistryKey)
}

func SetRegistry(s string) {
	setAndSave(viperProjectRegistryKey, s)
}

func Branch() string {
	return v.GetString(viperProjectBranchKey)
}

func SetBranch(s string) {
	setAndSave(viperProjectBranchKey, s)
}

func ApiVersion() string {
	return v.GetString(viperProjectApiVersionKey)
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
