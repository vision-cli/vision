package main

import (
	"os"

	"github.com/charmbracelet/log"

	"github.com/vision-cli/vision/plugins/vision-plugin-sample-v1/cmd/info"
	"github.com/vision-cli/vision/plugins/vision-plugin-sample-v1/cmd/version"
)

func main() {

	if len(os.Args) < 1 {
		log.Error("missing argument")
		return
	}

	switch os.Args[1] {
	case "":
		if err := cmd.InfoRootCmd.Execute(); err != nil {
			log.Error(err)
		}
	case "info":
		if err := cmd.InfoRootCmd.Execute(); err != nil {
			log.Error(err)
		}
	case "version":
		if err := cmd.VersionRootCmd.Execute(); err != nil {
			log.Error(err)
		}
	case "generate":

	}

	// switch os.Args[1] {
	// case "info":
	// 	json.NewEncoder(os.Stdout).Encode(map[string]string{
	// 		"short_description": "sample plugin short desc",
	// 	})
	// case "version":
	// 	json.NewEncoder(os.Stdout).Encode(map[string]string{
	// 		"sem_ver": "v0.0.1",
	// 	})
	// case "init":
	// 	json.NewEncoder(os.Stdout).Encode(map[string]any{
	// 		"config": map[string]string{
	// 			"subscriptionid": "value",
	// 		},
	// 	})
	// case "generate":
	// 	json.NewEncoder(os.Stdout).Encode("ACTION: generate new plugin")
	// }
}
