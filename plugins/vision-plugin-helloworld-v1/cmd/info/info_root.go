package cmd

import (
	_ "embed"
	"encoding/json"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	initialise "github.com/vision-cli/vision/plugins/vision-plugin-helloworld-v1/cmd/initialise"
)

func init() {
	initialise.InitRootCmd.PersistentFlags().VisitAll(visit)
}

type Flag struct {
	Name      string
	Shorthand string
	Usage     string
}

var flags []*Flag

var visit = func(f *pflag.Flag) {

	flags = append(flags, &Flag{
		Name:      f.Name,
		Shorthand: f.Shorthand,
		Usage:     f.Usage,
	})
	// log.Infof("visit all pflag function: %v", f)
}

var InfoRootCmd = &cobra.Command{
	Use:   "info",
	Short: "return info about the plugin",
	Long:  "ditto",
	RunE:  sampleCmd,
}

//go:embed info.txt
var infoOutput string

var sampleCmd = func(cmd *cobra.Command, args []string) error {

	json.NewEncoder(os.Stdout).Encode(map[string]any{
		"short_description": "a hello world example plugin",
		"long_description":  infoOutput,
		"flags":             flags,
	})

	return nil
}
