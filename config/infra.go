package config

const (
	viperInfraPrefix       = "infra"
	viperInfraDirectoryKey = viperInfraPrefix + ".directory"
)

func InfraDirectory() string {
	return v.GetString(viperInfraDirectoryKey)
}

func SetInfraDirectory(s string) {
	setAndSave(viperInfraDirectoryKey, s)
}
