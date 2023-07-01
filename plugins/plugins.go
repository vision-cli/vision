package plugins

import (
	"os"
	"os/exec"
	"strings"

	"github.com/vision-cli/vision/cli"
)

func GetPlugins() []string {
	var plugins []string
	pluginPath := goBinPath()
	pluginFiles, err := os.ReadDir(pluginPath)
	if err != nil {
		cli.Fatalf("Cannot read plugin directory %s", pluginPath)
	}
	for _, pluginFile := range pluginFiles {
		if !pluginFile.IsDir() && fileIsVisionPlugin(pluginFile.Name()) {
			plugins = append(plugins, pluginFile.Name())
		}
	}
	return plugins
}

func goBinPath() string {
	goBinPath := os.Getenv("GOBIN")
	if goBinPath == "" {
		goPath, err := exec.Command("go", "env", "GOPATH").Output()
		if err != nil {
			cli.Fatalf("Cannot determine GOBIN. Please set GOBIN or GOPATH")
		}
		goBinPath = string(goPath)[:len(goPath)-1] + "/bin"
	}
	return goBinPath
}

func fileIsVisionPlugin(filename string) bool {
	c := strings.Split(filename, "-")
	// eventually remove 'example'
	if len(c) != 4 || c[0] != "vision" || c[1] != "plugin" {
		return false
	}
	return true
}
