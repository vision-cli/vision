package config

import "github.com/spf13/viper"

type Persist interface {
	Set(key, value string)
	GetString(key string) string
	WriteConfig() error
	WriteConfigAs(filename string) error
	SetConfigName(name string)
	AddConfigPath(path string)
	ReadInConfig() error
}

type ViperPersist struct{}

func (v *ViperPersist) Set(key, value string) {
	viper.Set(key, value)
}

func (v *ViperPersist) GetString(key string) string {
	return viper.GetString(key)
}

func (v *ViperPersist) WriteConfig() error {
	return viper.WriteConfig()
}

func (v *ViperPersist) WriteConfigAs(filename string) error {
	return viper.WriteConfigAs(filename)
}

func (v *ViperPersist) SetConfigName(name string) {
	viper.SetConfigName(name)
}

func (v *ViperPersist) AddConfigPath(path string) {
	viper.AddConfigPath(path)
}

func (v *ViperPersist) ReadInConfig() error {
	return viper.ReadInConfig()
}
