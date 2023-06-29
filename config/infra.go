package config

import "github.com/spf13/viper"

const (
	viperInfraPrefix       = "infra"
	viperInfraDirectoryKey = viperInfraPrefix + ".directory"
)

func InfraDirectory() string {
	return viper.GetString(viperInfraDirectoryKey)
}

func SetInfraDirectory(s string) {
	setAndSave(viperInfraDirectoryKey, s)
}
