package upgrade

import (
	"fmt"

	"github.com/spf13/cobra"
)

func upgrade(cmd *cobra.Command, args []string) error {
	fmt.Println("this is the upgrade function")
	return nil
}

var UpgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "upgrade plugin",
	Long:  "upgrade plugin to latest version",
	RunE:  upgrade,
}
