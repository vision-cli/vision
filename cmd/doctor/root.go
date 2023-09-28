package doctor

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/vision-cli/vision/internal/plugin"
)

var RootCmd = &cobra.Command{
	Use:   "doctor",
	Short: "check status of plugins",
	Long:  "Check the status of all plugins and their commands. If a command fails, it will log an error.",
	RunE:  doctorCmd,
}

type ErrInvalidPlugin struct {
	PluginName string
	Command    string
	Reasons    []string
}

func (e ErrInvalidPlugin) Error() string {
	out := ""
	out += fmt.Sprintf("Plugin: %s\n", e.PluginName)
	out += fmt.Sprintf("Command: %s\n\n", e.Command)
	for _, reason := range e.Reasons {
		out += fmt.Sprintf("    ->  %s\n", reason)
	}
	return out
}

// doctorCmd looks through available plugins and checks for plugin health.
// If plugins commands are missing or incomplete, doctor returns them as faulty with a reason and prints them out.
var doctorCmd = func(cmd *cobra.Command, args []string) error {
	plugins := plugin.Find()

	var invalidPlugins []error

	for _, p := range plugins {
		// call each of the built in commands
		exe := plugin.NewExecutor(p.FullPath, args)
		var reasons []string
		info, err := exe.Info()
		// TODO(luke): add "not a string" catch to empty string checks
		switch {
		case err != nil:
			reasons = append(reasons, fmt.Sprintf("%v", err))
		case info.ShortDescription == "":
			reasons = append(reasons, "short description missing")
		case info.LongDescription == "":
			reasons = append(reasons, "log description missing")
		}

		if len(reasons) > 0 {
			invalidPlugins = append(invalidPlugins, ErrInvalidPlugin{
				PluginName: p.Name,
				Command:    "info",
				Reasons:    reasons,
			})
		}
	}
	for _, ip := range invalidPlugins {
		fmt.Println(ip.Error())
	}
	return nil
}
