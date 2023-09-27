package main

import (
	"encoding/json"
	"os"

	"github.com/charmbracelet/log"

	info "github.com/vision-cli/vision/plugins/vision-plugin-sample-v1/cmd/info"
	initialise "github.com/vision-cli/vision/plugins/vision-plugin-sample-v1/cmd/initialise"
	version "github.com/vision-cli/vision/plugins/vision-plugin-sample-v1/cmd/version"
)

func main() {

	if len(os.Args) < 1 {
		if err := info.InfoRootCmd.Execute(); err != nil {
			log.Error(err)
		}
		return
	}

	var rootCmd = "sample"

	switch os.Args[1] {
	case "":
		if err := info.InfoRootCmd.Execute(); err != nil {
			log.Error(err)
		}
	case "info":
		if err := info.InfoRootCmd.Execute(); err != nil {
			log.Error(err)
		}
	case "version":
		if err := version.VersionRootCmd.Execute(); err != nil {
			log.Error(err)
		}
	case "generate":
		json.NewEncoder(os.Stdout).Encode(map[string]bool{
			"success": true,
		})
	case "init":
		if err := initialise.InitRootCmd.Execute(); err != nil {
			log.Error(err)
		}
	case "help":
		// returns json
	}
}
