package upgrade

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/vision-cli/vision/plugins/vision-plugin-plugin-v1/config"
	"github.com/vision-cli/vision/plugins/vision-plugin-plugin-v1/internal"
)

func upgrade(cmd *cobra.Command, args []string) error {
	var config config.Config
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("reading in config: %w", err)
	}

	err := viper.Unmarshal(&config)
	if err != nil {
		return fmt.Errorf("unmarshalling config: %w", err)
	}

	exe := internal.Executor{
		PluginModule: config.PLUGIN_MODULE_URL,
	}

	// TODO(genevieve + luke): improve 29-35
	err = exe.UpgradeByGo()
	if err != nil {
		err = exe.UpgradeByCurl()
		if err != nil {
			fmt.Println("upgrading plugin: %w", err)
		}
	}

	return nil
}

var UpgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "upgrade plugin",
	Long:  "upgrade plugin to latest version",
	RunE:  upgrade,
}
