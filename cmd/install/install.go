package install

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"

	"github.com/charmbracelet/log"
)

var InstallCmd = &cobra.Command{
	Use:   "install [plugin]",
	Short: "Install a standard vision plugin",
	Long:  "Install one of the standard vision plugins: gorest, 10100 or plugin",
	RunE:  cmd,
}

var pluginMap = map[string]string{
	"gorest": "github.com/vision-cli/vision-plugin-gorest-v1@latest",
	"10100":  "github.com/vision-cli/vision-plugin-10100-v1@latest",
	"plugin": "github.com/vision-cli/vision-plugin-plugin-v1@latest",
}

var cmd = func(cmd *cobra.Command, args []string) error {
	for _, pluginName := range args {
		if _, ok := pluginMap[pluginName]; !ok {
			return fmt.Errorf("plugin %s is not a standard vision plugin. Please install it using go install <url>", pluginName)
		}
		err := install(pluginName)
		if err != nil {
			return err
		}
		log.Info(fmt.Sprintf("Installed plugin %s\n", pluginName))
	}
	return nil
}

func install(pluginName string) error {
	cmd := exec.Command("go", "install", pluginMap[pluginName])

	_, err := cmd.Output()
	if err != nil {
		return err
	}

	return nil
}
