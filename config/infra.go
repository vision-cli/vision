package config

const (
	viperInfraPrefix       = "infra"
	viperInfraDirectoryKey = viperInfraPrefix + ".directory"
)

func InfraDirectory() string {
	return v.GetStringOrDefault(viperInfraDirectoryKey, defaultInfraDir)
}

func SetInfraDirectory(s string) {
	setAndSave(viperInfraDirectoryKey, s)
}
