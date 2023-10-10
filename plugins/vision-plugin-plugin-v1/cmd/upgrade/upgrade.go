package upgrade

import (
	"github.com/spf13/cobra"

	"github.com/vision-cli/vision/plugins/vision-plugin-plugin-v1/internal"
)

func upgrade(cmd *cobra.Command, args []string) error {
	exe := internal.Executor{}

	exe.UpdateByCurl()
	return nil
}

var UpgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "upgrade plugin",
	Long:  "upgrade plugin to latest version",
	RunE:  upgrade,
}
