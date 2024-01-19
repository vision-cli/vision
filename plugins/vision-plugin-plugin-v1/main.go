package main

import (
	"github.com/spf13/viper"

	"github.com/vision-cli/vision/plugins/vision-plugin-plugin-v1/cmd"
)

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath("config/")
	viper.AutomaticEnv()
	viper.SetConfigType("yml")

	cmd.Execute()
}
